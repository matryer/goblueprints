package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results"`
}

func handlePolls(w http.ResponseWriter, r *http.Request) {

	db := GetVar(r, "db").(*mgo.Database)
	c := db.C("polls")

	switch r.Method {
	case "GET":

		var q *mgo.Query
		p := NewPath(r.URL.Path)
		if p.HasID() {
			// get specific poll
			q = c.FindId(bson.ObjectIdHex(p.ID))
		} else {
			// get all polls
			q = c.Find(nil)
		}
		var result []*poll
		if err := q.All(&result); err != nil {
			respondErr(w, r, err)
			return
		}
		respond(w, r, http.StatusOK, &result)
		return

	case "POST":

		var p poll
		if err := decodeBody(r, &p); err != nil {
			respondErr(w, r, "failed to read poll from request", err)
			return
		}
		p.ID = bson.NewObjectId()
		if err := c.Insert(p); err != nil {
			respondErr(w, r, "failed to insert poll", err)
			return
		}
		w.Header().Set("Location", "polls/"+p.ID.Hex())
		respond(w, r, http.StatusCreated, nil)
		return
	}

	// not found
	respondNotFound(w, r)

}
