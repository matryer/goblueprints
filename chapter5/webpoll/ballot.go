package webpoll

import "sync"

// Ballot represents a source of votes.
type Ballot interface {
	// Start should begin collecting votes in the background
	// and return a channel on which the votes will be reported.
	// If Start returns an error, it should close the channel itself.
	Start(options []string) (<-chan string, error)
	// Stop should stop the ballot from collecting votes.
	Stop()
}

// Ballots represents many Ballot objects.
type Ballots []Ballot

// Start starts all Ballots aggregating the results into a single
// channel which is returned.
// If any single Ballot fails to Start, the operation is aborted and
// the error is returned.
func (bs Ballots) Start(options []string) (<-chan string, error) {

	out := make(chan string)

	// keep track of started and finished ballots
	startErrs := make(map[int]error)
	var startX sync.Mutex
	var start, finish sync.WaitGroup
	start.Add(len(bs))
	finish.Add(len(bs))

	for i, ballot := range bs {
		go func(i int, ballot Ballot) {
			thisOut, err := ballot.Start(options)
			if err != nil {
				startX.Lock()
				startErrs[i] = err // save the error
				startX.Unlock()
				start.Done() // started but failed
			} else {
				start.Done() // started
				for vote := range thisOut {
					out <- vote
				}
			}
			finish.Done()
		}(i, ballot)
	}

	go func() {
		// wait for all ballots to finish
		finish.Wait()
		close(out)
	}()

	start.Wait() // wait for everything to have started

	// any errors?
	if len(startErrs) > 0 {
		// stop any ballots that successfully started
		var lastErr error
		for i, _ := range bs {
			if startErrs[i] == nil {
				bs[i].Stop()
			} else {
				lastErr = startErrs[i]
			}
		}
		return nil, lastErr
	}

	// no errors - return the channel
	return out, nil
}

// Stop stops all Ballots from collecting votes and causes
// the channel returned from Start to be closed.
func (bs Ballots) Stop() {
	for _, ballot := range bs {
		ballot.Stop()
	}
}

// Count counts the votes that come in on the channel.
// Blocks and returns a map[string]int of the results when the
// channel is closed.
func Count(votes <-chan string) map[string]int {
	counter := new(Counter)
	for _ = range counter.Count(votes) {
	}
	return counter.Results()
}

// Counter allows you to count the votes inline.
type Counter struct {
	results  map[string]int
	resultsX sync.Mutex
	out      chan string
}

// Count starts counting the votes from the channel and returns
// a new channel through which each vote will be sent.
func (c *Counter) Count(votes <-chan string) <-chan string {
	c.out = make(chan string)
	c.results = make(map[string]int)
	go func() {
		for vote := range votes {
			c.resultsX.Lock()
			c.results[vote]++ // increase the count
			c.resultsX.Unlock()
			c.out <- vote // send it out
		}
		c.Stop()
	}()
	return c.out
}

// Results gets the current results from the counting.
func (c *Counter) Results() map[string]int {
	c.resultsX.Lock()
	defer c.resultsX.Unlock()
	return c.results
}

// Stop stops counting.
func (c *Counter) Stop() {
	close(c.out)
}
