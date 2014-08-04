package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"labix.org/v2/mgo"

	"github.com/bitly/go-nsq"
	"github.com/joeshaw/envdecode"
	"github.com/matryer/go-oauth/oauth"
)

type twitterSettings struct {
	ConsumerKey    string `env:"WEBPOLL_TWITTER_KEY,required"`
	ConsumerSecret string `env:"WEBPOLL_TWITTER_SECRET,required"`
	AccessToken    string `env:"WEBPOLL_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"WEBPOLL_TWITTER_ACCESSSECRET,required"`
}

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
	conn       net.Conn
	reader     io.ReadCloser
	filterForm string
)

type poll struct {
	Options []string
}
type tweet struct {
	Text string `json:"text"`
}

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

func main() {

	var ts twitterSettings
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if conn != nil {
					conn.Close()
					conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				conn = netc
				return netc, nil
			},
		},
	}
	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}

	// handle stopping
	twitterStopChan := make(chan struct{}, 1)
	publisherStopChan := make(chan struct{}, 1)
	stop := false
	votes := make(chan string) // chan for votes

	// handle Ctrl+C
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stop = true
		log.Println("Stopping...")
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// publish votes on the "votes" channel
	go func() {
		pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // publish vote
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		publisherStopChan <- struct{}{}
	}()

	// read tweets
	go func() {
		defer func() {
			twitterStopChan <- struct{}{}
		}()
		for {

			if stop {
				log.Println("Twitter: Stopped")
				return
			}

			log.Println("Twitter: Waiting...")
			// wait so we don't overwhelm Twitter
			time.Sleep(2 * time.Second)
			log.Println("Twitter: Renewing")

			var options []string
			db, err := mgo.Dial("localhost")
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()
			iter := db.DB("ballots").C("polls").Find(nil).Iter()
			var p poll
			for iter.Next(&p) {
				options = append(options, p.Options...)
			}
			iter.Close()

			lowerOpts := make([]string, len(options))
			hashtags := make([]string, len(options))
			for i := range options {
				lowerOpts[i] = strings.ToLower(options[i])
				hashtags[i] = "#" + lowerOpts[i]
			}

			form := url.Values{"track": {strings.Join(hashtags, ",")}}
			formEnc := form.Encode()

			u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
			if err != nil {
				log.Println("creating filter request failed:", err)
			}
			req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", u, form))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error getting response:", err)
				continue
			}
			if resp.StatusCode != 200 {
				log.Println("StatusCode =", resp.StatusCode)
				continue
			}

			reader = resp.Body
			decoder := json.NewDecoder(reader)

			// keep reading
			for {
				var tweet *tweet
				if err := decoder.Decode(&tweet); err == nil {
					for i, option := range lowerOpts {
						if strings.Contains(strings.ToLower(tweet.Text), option) {
							votes <- options[i]
						}
					}
				} else {
					break
				}
			}

		}

	}()

	// update by forcing the connection to close
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			if stop {
				break
			}
		}
	}()

	// wait for everything to stop
	<-twitterStopChan // important to avoid panic
	close(votes)
	<-publisherStopChan
	log.Println("Everything: stopped")

}
