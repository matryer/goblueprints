package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {

	apiKey := os.Getenv("BHT_APIKEY")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {

		word := s.Text()
		response, err := http.Get("http://words.bighugelabs.com/api/2/" + apiKey + "/" + word + "/json")
		if err == nil {
			var data map[string]interface{}
			if json.NewDecoder(response.Body).Decode(&data) == nil {
				for _, valuesMap := range data {
					for _, values := range valuesMap.(map[string]interface{}) {
						for _, value := range values.([]interface{}) {
							fmt.Println(value)
						}
					}
				}
			}
		}

	}

}
