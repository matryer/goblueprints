package meander

import (
	"strings"
)

type j struct {
	Name       string
	PlaceTypes []string
}

func (j *j) Public() interface{} {
	return map[string]interface{}{
		"name":    j.Name,
		"journey": strings.Join(j.PlaceTypes, "|"),
	}
}

// Journeys represents the pre-set journeys data.
var Journeys = []interface{}{
	&j{Name: "a romantic day", PlaceTypes: []string{"park", "bar", "movie_theater", "restaurant", "florist", "taxi_stand"}},
	&j{Name: "a shopping spree", PlaceTypes: []string{"department_store", "cafe", "clothing_store", "jewelry_store", "shoe_store"}},
	&j{Name: "a night out", PlaceTypes: []string{"bar", "casino", "food", "bar", "night_club", "bar", "bar", "hospital"}},
	&j{Name: "a culture day", PlaceTypes: []string{"museum", "cafe", "cemetery", "library", "art_gallery"}},
	&j{Name: "a pamper day", PlaceTypes: []string{"hair_care", "beauty_salon", "cafe", "spa"}},
}
