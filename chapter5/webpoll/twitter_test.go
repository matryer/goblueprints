package webpoll

import (
	"log"
	"time"

	"testing"
)

func XTestTwitterBallot(t *testing.T) {

	ballots := Ballots{NewTwitterBallot("w7CAtRP9H8TbVsFMujy8K09CA", "94GO0Yi7Dl6kDSjDqOIu80fzcRESi1v6nz00ICBBl3JopsW58c", "29481227-WszI4ij0AZCfrQ7Gs3MVjdr1qsKJrWctMlc21x6ed", "czEnqnByEDC7nvlbo6G97Hym8KAkuSBc8ghWyTMthZOTV")}
	votes, err := ballots.Start([]string{"one", "two", "three"})
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		ballots.Stop()
	}()

	log.Println(Count(votes))

}
