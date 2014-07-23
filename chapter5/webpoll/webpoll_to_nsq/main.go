package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bitly/go-nsq"
	"github.com/joeshaw/envdecode"
	"github.com/matryer/goblueprints/chapter5/webpoll"
)

var (
	optionsList = flag.String("options", "", "Comma separated list of poll options")
	ballotsList = flag.String("ballots", "twitter", "Comma separated list of ballots")
	topic       = flag.String("topic", "votes", "nsq topic")
	nsqAddr     = flag.String("nsq", "127.0.0.1:4150", "nsq address")
)

func main() {

	flag.Parse()

	options := strings.Split(*optionsList, ",")
	if len(options) < 2 {
		fatal("Must provide at least two options.")
	}

	ballots := ballots()
	if len(ballots) == 0 {
		fatal("Need at least one source.")
	}

	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		ballots.Stop()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	votes, err := ballots.Start(options)
	if err != nil {
		fatal(err)
	}

	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(*nsqAddr, cfg)
	if err != nil {
		fatal(err)
	}

	for vote := range votes {
		log.Println(vote)
		producer.Publish(*topic, []byte(vote))
	}

	producer.Stop()

}

type twitterSettings struct {
	ConsumerKey    string `env:"WEBPOLL_TWITTER_KEY,required"`
	ConsumerSecret string `env:"WEBPOLL_TWITTER_SECRET,required"`
	AccessToken    string `env:"WEBPOLL_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"WEBPOLL_TWITTER_ACCESSSECRET,required"`
}

func ballots() webpoll.Ballots {

	var ballots webpoll.Ballots

	if strings.Contains(*ballotsList, "twitter") {

		var ts twitterSettings
		if err := envdecode.Decode(&ts); err != nil {
			fatal(err)
		}
		ballots = append(ballots, webpoll.NewTwitterBallot(
			ts.ConsumerKey,
			ts.ConsumerSecret,
			ts.AccessToken,
			ts.AccessSecret,
		))

	}

	return ballots
}

func fatal(a ...interface{}) {
	fmt.Println(a...)
	flag.PrintDefaults()
	os.Exit(1)
}
