package api

import (
	"errors"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func handleVotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	handleVote(w, r)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var newVote struct {
		AnswerID string `json:"answer_id"`
		Score    int    `json:"score"`
	}
	err := decode(r, &newVote)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}
	err = validScore(newVote.Score)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}
	answerKey, err := datastore.DecodeKey(newVote.AnswerID)
	if err != nil {
		respondErr(ctx, w, r, errors.New("invalid answer_id"), http.StatusBadRequest)
		return
	}
	vote, err := CastVote(ctx, answerKey, newVote.Score)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, vote, http.StatusCreated)
}
