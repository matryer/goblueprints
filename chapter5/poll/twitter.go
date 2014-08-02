package poll

import (
	"encoding/json"
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

// TwitterBallot represents a ballot that pulls votes from
// Twitter.
type TwitterBallot struct {
	authClient *oauth.Client
	creds      *oauth.Credentials
	conn       net.Conn
	reader     io.ReadCloser
	out        chan string
}

// NewTwitterBallot makes a new Ballot capable of pulling requests
// from Twitter.
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

// Start starts the Ballot.
func (t *TwitterBallot) Start(options []string) (<-chan string, error) {

	lowerOpts := make([]string, len(options))
	hashtags := make([]string, len(options))
	for i := range options {
		lowerOpts[i] = strings.ToLower(options[i])
		hashtags[i] = "#" + lowerOpts[i]
	}

	form := url.Values{"track": {strings.Join(hashtags, ",")}}
	formEnc := form.Encode()
	u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if t.conn != nil {
					t.conn.Close()
					t.conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				t.conn = netc
				return netc, nil
			},
		},
	}

	t.out = make(chan string)

	/*
		// open the first connection
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("filter failed: %d", resp.StatusCode)
		}
	*/

	first := true
	go func() {
		for {

			if !first {
				// give the service some time to recover
				log.Println("Waiting (so we don't bother Twitter too much)")
				time.Sleep(10 * time.Second)
				log.Println("  carrying on...")
			} else {
				first = false
			}

			// if we have no response - make a new connection

			req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
			if err != nil {
				log.Println("creating filter request failed:", err)
			}
			req.Header.Set("Authorization", t.authClient.AuthorizationHeader(t.creds, "POST", u, form))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error getting response:", err)
				continue
			}
			if resp.StatusCode != 200 {
				log.Println("StatusCode==", resp.StatusCode)
				continue
			}

			t.reader = resp.Body
			decoder := json.NewDecoder(t.reader)

			// keep reading
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

		}
	}()

	return t.out, nil
}

// Stop stops the ballot.
func (t *TwitterBallot) Stop() {
	if t.conn == nil {
		return
	}
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
