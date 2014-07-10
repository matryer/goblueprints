package thesaurus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type BigHugh struct {
	APIKey string
}

func (b *BigHugh) Synonyms(term string) ([]string, error) {
	var syns []string
	response, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKey + "/" + term + "/json")
	if err != nil {
		return nil, errors.New("bighugh: Failed when looking for synonyms for \"" + term + "\"" + err.Error())
	}
	var data map[string]interface{}
	if json.NewDecoder(response.Body).Decode(&data) == nil {
		for _, valuesMap := range data {
			for key, values := range valuesMap.(map[string]interface{}) {
				if key == "syn" {
					for _, value := range values.([]interface{}) {
						syns = append(syns, value.(string))
					}
				}
			}
		}
	}
	return syns, nil
}
