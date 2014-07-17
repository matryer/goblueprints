package webpoll

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/matryer/go-oauth/oauth"
)

type tweet struct {
	Text string `json:"text"`
}

type TwitterBallot struct {
	authClient *oauth.Client
	creds      *oauth.Credentials
	conn       net.Conn
	reader     io.ReadCloser
	out        chan string
}

func NewTwitterBallot(consumerKey, consumerSecret, accessToken, accessSecret string) *TwitterBallot {
	return &TwitterBallot{
		creds: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessSecret,
		},
		authClient: &oauth.Client{
			Credentials: oauth.Credentials{
				Token:  consumerKey,
				Secret: consumerSecret,
			},
		},
	}
}

func (t *TwitterBallot) Start(options []string) (<-chan string, error) {

	lowerOpts := make([]string, len(options))
	hashtags := make([]string, len(options))
	for i, _ := range options {
		lowerOpts[i] = strings.ToLower(options[i])
		hashtags[i] = "#" + lowerOpts[i]
	}

	form := url.Values{"track": {strings.Join(hashtags, ",")}}
	formEnc := form.Encode()
	u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
	if err != nil {
		return nil, fmt.Errorf("Creating filter request failed: %s", err)
	}

	req.Header.Set("Authorization", t.authClient.AuthorizationHeader(t.creds, "POST", u, form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if t.conn != nil {
					t.conn.Close()
					t.conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 1*time.Minute)
				if err != nil {
					return nil, err
				}
				t.conn = netc
				return netc, nil
			},
		},
	}

	t.out = make(chan string)

	// open the first connection
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("Twitter filter failed: %d", resp.StatusCode)
	}

	go func() {
		for {

			// if we have no response - make a new connection
			if resp == nil {
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				if resp.StatusCode != 200 {
					continue
				}
			}

			t.reader = resp.Body
			decoder := json.NewDecoder(t.reader)

			for {
				var tweet *tweet
				if err := decoder.Decode(&tweet); err == nil {
					for i, option := range lowerOpts {
						if strings.Contains(strings.ToLower(tweet.Text), option) {
							t.out <- options[i]
						}
					}
				} else {
					// connection probably died - reconnect
					resp = nil
					break
				}
			}

			// give the service some time to recover
			time.Sleep(30 * time.Second)

		}
	}()

	return t.out, nil
}

func (t *TwitterBallot) Stop() {
	if err := t.conn.Close(); err != nil {
		log.Println("ERROR: Failed to close conn:", err)
	}
	if t.reader != nil {
		if err := t.reader.Close(); err != nil {
			log.Println("ERROR: Failed to close connection reader:", err)
		}
	}
	close(t.out)
}
