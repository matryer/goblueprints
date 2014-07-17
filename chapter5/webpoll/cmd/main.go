package main

import (
	"fmt"
	"os"
	"strings"

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
	options = flag.String("options", "", "Comma separated list of poll options")
	sources = flag.String("sources", "twitter", "Comma separated list of ballot sources")
)

type twitterSettings struct {
	ConsumerKey    string `env:"WEBPOLL_TWITTER_KEY,required"`
	ConsumerSecret string `env:"WEBPOLL_TWITTER_SECRET,required"`
	AccessToken    string `env:"WEBPOLL_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"WEBPOLL_TWITTER_ACCESSSECRET,required"`
}

func main() {

	flag.Parse()

	options := strings.Split(*options, ",")
	if len(options) < 2 {
		fatal("Must provide at least two options.")
	}

	var ballots webpoll.Ballots

	if strings.Contains(*sources, "twitter") {

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

}

func fatal(a ...interface{}) {
	fmt.Println(a...)
	flag.PrintDefaults()
	os.Exit(1)
}
