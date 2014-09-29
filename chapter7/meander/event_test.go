package meander_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/goblueprints/chapter7/meander"
	"github.com/stretchr/testify/require"
)

func TestEventJSON(t *testing.T) {

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

	require.NoError(t, json.NewDecoder(strings.NewReader(j)).Decode(&e))

	require.Equal(t, "Rhythmboat Cruises", e.Name)
	require.Equal(t, "http://maps.gstatic.com/mapfiles/place_api/icons/travel_agent-71.png", e.Icon)
	require.Equal(t, -33.870775, e.Lat)
	require.Equal(t, 151.199025, e.Lng)

}
