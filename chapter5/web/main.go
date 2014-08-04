package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"

	"strings"
	"time"

	"github.com/stretchr/graceful"

	"labix.org/v2/mgo"
)

var host = flag.String("host", ":8080", "The host of the application.")

var (
	dbSession *mgo.Session
)

func WithData(fn func(*mgo.Database, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thisDb := dbSession.Copy()
		fn(thisDb.DB("ballots"), w, r)
	})
}

func main() {

	// open database connection
	var err error
	dbSession, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatalln("Couldn't connect to database", err)
	}
	defer dbSession.Close()

	server := http.NewServeMux()
	server.HandleFunc("/polls.json", WithData(func(db *mgo.Database, w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET": // read polls
			var polls []map[string]interface{}
			err := db.C("polls").Find(nil).Sort("-_id").All(&polls)
			if err != nil {
				w.WriteHeader(500)
				io.WriteString(w, "reading: "+err.Error())
				return
			}
			json.NewEncoder(w).Encode(polls)
		case "POST": // create poll
			title := r.FormValue("title")
			options := strings.Fields(r.FormValue("options"))
			err := db.C("polls").Insert(map[string]interface{}{"title": title, "options": options})
			if err != nil {
				w.WriteHeader(500)
				io.WriteString(w, "creating: "+err.Error())
				return
			}
			w.Header()["Location"] = []string{"/"}
			w.WriteHeader(301)
		}
	}))
	server.HandleFunc("/polls/", WithData(func(db *mgo.Database, w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			segs := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			err := db.C("polls").RemoveId(bson.ObjectIdHex(segs[1]))
			if err != nil {
				w.WriteHeader(500)
				io.WriteString(w, "removing: "+err.Error())
				return
			}
			w.WriteHeader(200)
		}
	}))
	server.Handle("/",
		http.StripPrefix("/",
			http.FileServer(http.Dir("./files"))))
	log.Println("Starting web server on", *host)
	graceful.Run(*host, 10*time.Second, server)

	log.Println("Goodbye")

}
