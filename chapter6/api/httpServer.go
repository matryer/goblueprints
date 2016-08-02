package main

import (
	"net/http"

	"golang.org/x/net/context"
)

// httpServer is an implementation for storing polls
// via http with an agnostic storage backend
type httpServer struct {
	storage PollStorage
}

func NewHttpServer(storage PollStorage) *httpServer {
	return &httpServer{storage: storage}
}

func (s *httpServer) Routes(w http.ResponseWriter, r *http.Request) {
	withCORS(withAPIKey(s.handlePolls))
}

func (s *httpServer) handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlePollsGet(w, r)
		return
	case "POST":
		s.handlePollsPost(w, r)
		return
	case "DELETE":
		s.handlePollsDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *httpServer) handlePollsGet(w http.ResponseWriter, r *http.Request) {
	p := NewPath(r.URL.Path)
	result, err := s.storage.Get(p)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}

func (s *httpServer) handlePollsPost(w http.ResponseWriter, r *http.Request) {
	var p Poll
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read poll from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = apikey
	}
	id, err := s.storage.Create(&p)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert poll", err)
		return
	}
	w.Header().Set("Location", "polls/"+id)
	respond(w, r, http.StatusCreated, nil)
}

func (s *httpServer) handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	p := NewPath(r.URL.Path)

	err := s.storage.Delete(p)

	if err == ErrCannotDeleteAll {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all polls.")
		return
	}

	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}
	respond(w, r, http.StatusOK, nil) // ok
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

func APIKey(ctx context.Context) (string, bool) {
	key := ctx.Value(contextKeyAPIKey)
	if key == nil {
		return "", false
	}
	keystr, ok := key.(string)
	return keystr, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}
