package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/matryer/goblueprints/chapter5/webpoll"

	"github.com/joeshaw/envdecode"

	"flag"
)

/*
  webpoll command

  Usage:

    webpoll -options="one,two,three,four"

*/

var (
	optionsFlg = flag.String("options", "", "Comma separated list of poll options")
	sourcesFlg = flag.String("sources", "twitter", "Comma separated list of ballot sources")
	outputFlg  = flag.String("out", "./votes.json", "Output file for results")
	refresh    = flag.Int64("refresh", 10, "Seconds between updating the out results file")
)

func main() {

	flag.Parse()

	options := strings.Split(*optionsFlg, ",")
	if len(options) < 2 {
		fatal("Must provide at least two options.")
	}

	ballots := ballots()
	if len(ballots) == 0 {
		fatal("Need at least one source.")
	}

	votes, err := ballots.Start(options)
	if err != nil {
		fatal(err)
	}
	counter := new(webpoll.Counter)

	go func() {
		for vote := range counter.Count(votes) {
			fmt.Print(vote + " ")
		}
	}()

	// keep updating the results
	for {

		// wait
		time.Sleep(time.Duration(*refresh) * time.Second)

		// write the outputs
		results := counter.Results()
		resultsJson, _ := json.Marshal(results)
		if outFile, err := os.Create(*outputFlg); err == nil {
			outFile.Write(resultsJson)
			outFile.Close()
		} else {
			fmt.Println("couldn't open", *outputFlg, "-", err)
		}

	}

}

type twitterSettings struct {
	ConsumerKey    string `env:"WEBPOLL_TWITTER_KEY,required"`
	ConsumerSecret string `env:"WEBPOLL_TWITTER_SECRET,required"`
	AccessToken    string `env:"WEBPOLL_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"WEBPOLL_TWITTER_ACCESSSECRET,required"`
}

func ballots() webpoll.Ballots {

	var ballots webpoll.Ballots

	if strings.Contains(*sourcesFlg, "twitter") {

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
