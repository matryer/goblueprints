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
	results := make(map[string]int)
	for vote := range votes {
		results[vote]++
	}
	return results
}
