package main

/*

  poll command

  Connects to local mongo and looks for ballots.polls
  and publishes all realtime votes from Twitter on the
  "votes" nsq channel.

  Every minute, looks for new polls.

*/

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"labix.org/v2/mgo"

	"github.com/bitly/go-nsq"
	"github.com/joeshaw/envdecode"
	"github.com/matryer/goblueprints/chapter5/poll"

	"flag"
	"fmt"
)

type twitterSettings struct {
	ConsumerKey    string `env:"WEBPOLL_TWITTER_KEY,required"`
	ConsumerSecret string `env:"WEBPOLL_TWITTER_SECRET,required"`
	AccessToken    string `env:"WEBPOLL_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"WEBPOLL_TWITTER_ACCESSSECRET,required"`
}

type pollItem struct {
	Options []string
}

var stop bool
var stopLock sync.Mutex

func main() {

	// add twitter ballot
	var ts twitterSettings
	if err := envdecode.Decode(&ts); err != nil {
		fatal(err)
	}
	ballot := poll.NewTwitterBallot(
		ts.ConsumerKey,
		ts.ConsumerSecret,
		ts.AccessToken,
		ts.AccessSecret,
	)

	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		log.Println("Stopping...")
		stopLock.Lock()
		stop = true
		stopLock.Unlock()
		ballot.Stop()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// connect to the database
	log.Println("Connecting to db...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
	}
	log.Println("Connected")
	defer db.Close()
	pollData := db.DB("ballots").C("polls")

	// make a producer to send messages to nsq
	pub, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())

	for {

		stopLock.Lock()
		if stop {
			stopLock.Unlock()
			break
		}
		stopLock.Unlock()

		// get all options we care about from the
		// polls in the database
		var p pollItem
		var options []string
		pollCount, _ := pollData.Find(nil).Count()
		if pollCount == 0 {
			log.Println("No polls - will try again in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		} else {
			log.Println("Found", pollCount, "poll(s)...")
		}
		iter := pollData.Find(nil).Iter()
		for iter.Next(&p) {
			options = append(options, p.Options...)
		}
		iter.Close()

		// stop the ballot in 1 minute... it will then
		// restart
		go func() {
			time.Sleep(1 * time.Minute)
			log.Println("Restarting...")
			ballot.Stop()
		}()

		log.Println("Looking for:", strings.Join(options, ", "))

		votes, err := ballot.Start(options)
		if err != nil {
			fatal("ballot failed to start:", err)
		}

		log.Println("Waiting for votes to push to nsq...")

		for vote := range votes {
			fmt.Print(vote + " ")
			pub.Publish("votes", []byte(vote)) // publish vote
		}

		time.Sleep(10 * time.Second) // wait a while so we don't overwhelm Twitter

	}

	pub.Stop() // close nsq
	log.Println("Stopped")

}

func fatal(a ...interface{}) {
	fmt.Println(a...)
	flag.PrintDefaults()
	os.Exit(1)
}
