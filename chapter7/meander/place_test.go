package meander_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/goblueprints/chapter7/meander"
)

func TestPlaceJSON(t *testing.T) {
	is := is.New(t)

	var e meander.Place

	j := `{
   "geometry" : {
      "location" : {
         "lat" : -33.870775,
         "lng" : 151.199025
      }
   },
   "icon" : "http://maps.gstatic.com/mapfiles/place_api/icons/travel_agent-71.png",
   "id" : "21a0b251c9b8392186142c798263e289fe45b4aa",
   "name" : "Rhythmboat Cruises"
  }`

	is.NoErr(json.NewDecoder(strings.NewReader(j)).Decode(&e))

	is.Equal("Rhythmboat Cruises", e.Name)
	is.Equal("http://maps.gstatic.com/mapfiles/place_api/icons/travel_agent-71.png", e.Icon)
	is.Equal(-33.870775, e.Lat)
	is.Equal(151.199025, e.Lng)

}
