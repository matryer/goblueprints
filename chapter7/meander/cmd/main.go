package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/matryer/goblueprints/chapter7/meander"
)

func main() {
	meander.APIKey = "AIzaSyBw8jI_Lq_UjCyHqKaxtoyYcNVGIeJG1fE"
	http.HandleFunc("/places", func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{
			Route: strings.Split(r.URL.Query().Get("route"), "|"),
		}
		q.Lat, _ = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		q.Lng, _ = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		q.Radius, _ = strconv.Atoi(r.URL.Query().Get("radius"))
		q.CostRangeStr = r.URL.Query().Get("cost")
		places := q.Run()
		json.NewEncoder(w).Encode(places)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
