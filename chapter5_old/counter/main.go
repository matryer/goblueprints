package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"
	"labix.org/v2/mgo"
)

var (
	topic    = flag.String("topic", "votes", "nsq topic")
	nsqAddr  = flag.String("nsq", "127.0.0.1:4150", "nsq address")
	interval = flag.Int("interval", 1, "(seconds) update interval")
)

var (
	results  map[string]int
	resultsM sync.Mutex
	pollsCol *mgo.Collection
)

func main() {

	flag.Parse()
	if len(*topic) == 0 {
		fatal("Needs a valid topic")
	}

	var fatalErr error
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(*topic, "counter", cfg)
	if err != nil {
		fatal(err)
	}
	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		vote := string(m.Body)
		resultsM.Lock()
		if results == nil {
			results = make(map[string]int)
		}
		results[vote]++
		resultsM.Unlock()
		return nil
	}))
	if err := consumer.ConnectToNSQD(*nsqAddr); err != nil {
		fatal(err)
	} else {
		fmt.Println("Connected")
	}

	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
	}
	defer db.Close()

	clock := time.NewTicker(time.Duration(*interval) * time.Second)
	pollsCol = db.DB("webpoll").C("polls")
	go func() {
		for _ = range clock.C {
			update()
		}
	}()

	// wait for things to finish
	<-termChan
	clock.Stop()
	consumer.Stop()
	<-consumer.StopChan

	// if a fatal error occurred - report it
	if fatalErr != nil {
		fatal(fatalErr)
	}

}

// update saves the latest results to the database
func update() {
	resultsM.Lock()
	up := map[string]interface{}{"$inc": results}
	if _, err := pollsCol.UpsertId(*topic, up); err != nil {
		fmt.Println("Error saving update to db:", err)
	} else {
		fmt.Println("Results saved:", results)
	}
	results = nil
	resultsM.Unlock()
}

func fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}
