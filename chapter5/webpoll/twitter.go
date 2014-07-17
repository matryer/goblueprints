package webpoll

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	for i, _ := range options {
		options[i] = strings.ToLower(options[i])
	}

	form := url.Values{"track": {strings.Join(options, ",")}}
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
				netc, err := net.DialTimeout(netw, addr, 1*time.Minute)
				if err != nil {
					return nil, err
				}
				t.conn = netc
				return netc, nil
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Twitter filter failed (%d): %s", resp.StatusCode, body)
	}

	t.reader = resp.Body
	decoder := json.NewDecoder(t.reader)

	t.out = make(chan string)

	go func() {
		for {
			var tweet *tweet
			if err := decoder.Decode(&tweet); err != nil {
				return
			} else {
				for _, option := range options {
					if strings.Contains(strings.ToLower(tweet.Text), option) {
						t.out <- option
					}
				}
			}
		}
	}()

	return t.out, nil
}

func (t *TwitterBallot) Stop() {
	if t.reader != nil {
		t.reader.Close()
	}
	t.conn.Close()
	close(t.out)
}
