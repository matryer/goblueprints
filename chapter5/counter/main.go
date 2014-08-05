package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
  counter tool looks for votes on the "votes" nsq topic
  and adds the numbers to the polls in mongo
*/

const updateDuration = 1 * time.Second

type pollItem struct {
	Options map[string]int
}

func main() {

	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// connect to the database
	log.Println("Connecting to database...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("Closing database connection...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	// listen for votes
	log.Println("Connecting to nsq...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	var counts map[string]int
	var countsLock sync.Mutex

	// when we see a vote...
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		q.Stop()
		<-q.StopChan
		return
	}

	log.Println("Waiting for votes on nsq...")

	// push to the database every updateDuration
	var updater *time.Timer
	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("No new votes, skipping database update")
		} else {
			fmt.Println()
			log.Println("Updating database...")
			log.Println(counts)
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if err := pollData.Update(sel, up); err != nil {
					log.Println("failed to update:", err)
				}
			}
			log.Println("Finished updating database...")
			counts = nil // reset counts
		}
		updater.Reset(updateDuration)
	})

	for {
		select {
		case <-termChan:
			updater.Stop()
			q.Stop()
		case <-q.StopChan:
			// finished
			return
		}
	}

}

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}
