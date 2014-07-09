package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		counter := 0
		response, err := http.Get("http://words.bighugelabs.com/api/2/" + apiKey + "/" + word + "/json")
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for \""+word+"\"", err)
		}
		var data map[string]interface{}
		if json.NewDecoder(response.Body).Decode(&data) == nil {
			for _, valuesMap := range data {
				for _, values := range valuesMap.(map[string]interface{}) {
					for _, value := range values.([]interface{}) {
						fmt.Println(value)
						counter++
					}
				}
			}
		}
		if counter == 0 {
			log.Fatalln("Couldn't find any synonyms for \"" + word + "\"")
		}
	}
}
