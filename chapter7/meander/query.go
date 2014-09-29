package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var APIKey string

type Query struct {
	Lat          float64
	Lng          float64
	Route        []string
	Radius       int
	CostRangeStr string
}

func (q *Query) find(types string) (*googleResponse, error) {
	u := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	u = fmt.Sprintf("%s?location=%g,%g&radius=%d", u, q.Lat, q.Lng, q.Radius)
	u = fmt.Sprintf("%s&types=%s", u, types)
	u = fmt.Sprintf("%s&key=%s", u, APIKey)
	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		u = fmt.Sprintf("%s&minprice=%d&maxprice=%d", u, int(r.From), int(r.To))
	}
	log.Println(u)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var response googleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (q *Query) Run() []*Place {
	rand.Seed(time.Now().UnixNano())
	var w sync.WaitGroup
	var eventsL sync.Mutex
	events := make([]*Place, len(q.Route))
	for i, r := range q.Route {
		w.Add(1)
		go func(types string, i int) {
			defer w.Done()
			response, err := q.find(types)
			if err != nil {
				log.Println("Failed to find places:", err)
				return
			}
			if len(response.Results) == 0 {
				log.Println("No places found for", types)
				return
			}
			for _, result := range response.Results {
				for _, photo := range result.Photos {
					photo.URL = "https://maps.googleapis.com/maps/api/place/photo?" +
						"maxwidth=400&photoreference=" + photo.PhotoRef + "&key=" + APIKey
				}
			}
			eventsL.Lock()
			defer eventsL.Unlock()
			randI := rand.Intn(len(response.Results))
			events[i] = response.Results[randI]
		}(r, i)
	}
	w.Wait() // wait for everything to finish
	return events
}

type googleResponse struct {
	Results []*Place `json:"results"`
}
type Place struct {
	Name            string `json:"name"`
	Icon            string `json:"icon"`
	*googleGeometry `json:"geometry"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}
type googleGeometry struct {
	*googleLocation `json:"location"`
}
type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type googlePhoto struct {
	Height   int    `json:"height"`
	Width    int    `json:"Width"`
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}
