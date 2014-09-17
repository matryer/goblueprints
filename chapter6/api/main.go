package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/stretchr/graceful"
	"gopkg.in/mgo.v2"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Starting web server on", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("Stopping...")
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		fn(w, r)
	}
}

func withData(db *mgo.Session, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := db.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}
