package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/stretchr/graceful"
	"labix.org/v2/mgo"
)

var (
	dbSession *mgo.Session
)

func WithData(fn func(*mgo.Database, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thisDb := dbSession.Copy()
		fn(thisDb.DB("webpoll"), w, r)
	})
}

func main() {

	var err error
	dbSession, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatalln("Couldn't connect to database", err)
	}
	defer dbSession.Close()

	var templ *template.Template
	templ = template.Must(template.ParseFiles("templates/index.html"))

	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, r)
	})
	s.HandleFunc("/polls.json", WithData(func(db *mgo.Database, w http.ResponseWriter, r *http.Request) {
		var polls []map[string]interface{}
		err := db.C("polls").Find(nil).Sort("_id").All(&polls)
		if err != nil {
			io.WriteString(w, err.Error())
		}
		json.NewEncoder(w).Encode(polls)
	}))
	graceful.Run(":8080", 10*time.Second, s)

}
