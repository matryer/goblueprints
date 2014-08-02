package poll_test

import (
	"errors"
	"testing"
	"time"

	"github.com/matryer/goblueprints/chapter5/poll"
)

var ErrTest = errors.New("test error")

type TestBallot struct {
	out       chan string
	ShouldErr bool
}

// TestBallot and poll.Ballots must implement poll.Ballot
var _ poll.Ballot = (*TestBallot)(nil)
var _ poll.Ballot = (poll.Ballots)(nil)

func (b *TestBallot) Start(options []string) (<-chan string, error) {
	b.out = make(chan string)

	if b.ShouldErr {
		close(b.out)
		return b.out, ErrTest
	}

	return b.out, nil
}
func (b *TestBallot) Stop() {
	close(b.out)
}
func (b *TestBallot) Vote(option string) {
	select {
	case b.out <- option:
	}
}

func TestBallotImpl(t *testing.T) {

	b := &TestBallot{}
	out, _ := b.Start([]string{"one", "two", "three"})
	if out == nil {
		t.Error("Ballots.Start should return a channel")
	}

	// stop after 1/2 second
	go func() {
		time.Sleep(500 * time.Millisecond)
		b.Stop()
	}()

	options := []string{"one", "two", "three"}

	// simulate real data
	go func() {
		for _, v := range options {
			b.Vote(v)
		}
	}()

	var actualOptions []string
	for option := range out {
		actualOptions = append(actualOptions, option)
	}

	for i := range options {
		if options[i] != actualOptions[i] {
			t.Error("Option", i, "should be", options[i], "but was", actualOptions[i])
		}
	}

}

func TestManyBallotStartErr(t *testing.T) {

	b1 := &TestBallot{}
	b2 := &TestBallot{}
	b3 := &TestBallot{}

	b2.ShouldErr = true

	bs := poll.Ballots([]poll.Ballot{b1, b2, b3})
	out, err := bs.Start([]string{"one", "two", "three"})

	if err != ErrTest {
		t.Error("Expected ErrTest error")
	}
	if out != nil {
		t.Error("No channel expected if an error occurred")
	}

}

func TestManyBallot(t *testing.T) {

	b1 := &TestBallot{}
	b2 := &TestBallot{}
	b3 := &TestBallot{}

	bs := poll.Ballots([]poll.Ballot{b1, b2, b3})
	out, _ := bs.Start([]string{"one", "two", "three"})

	if out == nil {
		t.Error("Ballots.Start should return a channel")
	}

	// stop after 1/2 second
	go func() {
		time.Sleep(500 * time.Millisecond)
		bs.Stop()
	}()

	options := []string{"one", "two", "three"}
	// simulate real data
	go func() {
		for _, v := range options {
			b1.Vote(v)
			b2.Vote(v)
			b3.Vote(v)
		}
	}()

	voteCount := make(map[string]int)
	for option := range out {
		voteCount[option]++
	}

	if voteCount["one"] != 3 {
		t.Error("Expected 3 x one")
	}
	if voteCount["two"] != 3 {
		t.Error("Expected 3 x two")
	}
	if voteCount["three"] != 3 {
		t.Error("Expected 3 x three")
	}

}
