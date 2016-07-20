package meander

import (
	"strings"
)

// Journey represents a journey template.
type Journey struct {
	Name       string
	PlaceTypes []string
}

// Public gets a public view of this Journey.
func (j Journey) Public() interface{} {
	return map[string]interface{}{
		"name":    j.Name,
		"journey": strings.Join(j.PlaceTypes, "|"),
	}
}

// Journeys represents the pre-set journeys data.
var Journeys = []interface{}{
	Journey{Name: "a romantic day", PlaceTypes: []string{"park", "bar", "movie_theatre", "restaurant", "florist", "taxi_stand"}},
	Journey{Name: "a shopping spree", PlaceTypes: []string{"department_store", "cafe", "clothing_store", "jewelry_store", "shoe_store"}},
	Journey{Name: "a night out", PlaceTypes: []string{"bar", "casino", "food", "bar", "night_club", "bar", "bar", "hospital"}},
	Journey{Name: "a culture day", PlaceTypes: []string{"museum", "cafe", "cemetery", "library", "art_gallery"}},
	Journey{Name: "a pamper day", PlaceTypes: []string{"hair_care", "beauty_salon", "cafe", "spa"}},
}
