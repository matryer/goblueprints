package api

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleQuestionCreate(w, r)
	case "GET":
		handleQuestionGet(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleQuestionGet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	params := pathParams(r, "/api/questions/:id")
	questionID, ok := params[":id"]
	if !ok {
		handleTopQuestions(w, r)
		return
	}
	questionKey, err := datastore.DecodeKey(questionID)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}
	question, err := GetQuestion(ctx, questionKey)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			respondErr(ctx, w, r, datastore.ErrNoSuchEntity, http.StatusNotFound)
			return
		}
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, question, http.StatusOK)
}

func handleTopQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	questions, err := TopQuestions(ctx)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, questions, http.StatusOK)
}

func handleQuestionCreate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var q Question
	err := decode(r, &q)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}
	err = q.Put(ctx)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, q, http.StatusCreated)
}
