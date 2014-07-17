package webpoll_test

import (
	"errors"
	"testing"
	"time"

	"github.com/matryer/goblueprints/chapter5/webpoll"
)

var ErrTest = errors.New("This is a test error")

type TestBallot struct {
	out       chan string
	ShouldErr bool
}

// TestBallot and webpoll.Ballots must implement webpoll.Ballot
var _ webpoll.Ballot = (*TestBallot)(nil)
var _ webpoll.Ballot = (webpoll.Ballots)(nil)

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
	b.out <- option
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

	for i, _ := range options {
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

	bs := webpoll.Ballots([]webpoll.Ballot{b1, b2, b3})
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

	bs := webpoll.Ballots([]webpoll.Ballot{b1, b2, b3})
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

func TestCount(t *testing.T) {

	b1 := &TestBallot{}
	b2 := &TestBallot{}
	b3 := &TestBallot{}

	bs := webpoll.Ballots([]webpoll.Ballot{b1, b2, b3})
	out, _ := bs.Start([]string{"two", "three", "five", "none"})

	// stop after 1/2 second
	go func() {
		time.Sleep(500 * time.Millisecond)
		bs.Stop()
	}()

	// simulate real concurrent data
	go func() {
		b1.Vote("three")
		b3.Vote("five")
		b1.Vote("two")
		b1.Vote("five")
	}()
	go func() {
		b2.Vote("three")
		b3.Vote("five")
		b2.Vote("five")
	}()
	go func() {
		b3.Vote("three")
		b3.Vote("two")
		b3.Vote("five")
	}()

	results := webpoll.Count(out)

	if results["three"] != 3 {
		t.Error("Expected 3x'three', but there was", results["three"])
	}
	if results["two"] != 2 {
		t.Error("Expected 2x'two', but there was", results["two"])
	}
	if results["five"] != 5 {
		t.Error("Expected 5x'five', but there was", results["five"])
	}
	if results["none"] != 0 {
		t.Error("Expected 0x'none', but there was", results["five"])
	}

}
