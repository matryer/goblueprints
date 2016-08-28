package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

// decode decodes the request body into v and calls
// its OK() function to validate the object.
func decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	if valid, ok := v.(interface {
		OK() error
	}); ok {
		err = valid.OK()
		if err != nil {
			return err
		}
	}
	return nil
}

func respond(ctx context.Context, w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Errorf(ctx, "respond: %s", err)
	}
}

func respondErr(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, code int) {
	errObj := struct {
		Error string `json:"error"`
	}{Error: err.Error()}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		log.Errorf(ctx, "respondErr: %s", err)
	}
}

// pathParams parses the URL.Path of the Request with the given
// pattern, and extracts the value for each segment into a map.
func pathParams(r *http.Request, pattern string) map[string]string {
	params := map[string]string{}
	pathSegs := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, seg := range strings.Split(strings.Trim(pattern, "/"), "/") {
		if i > len(pathSegs)-1 {
			return params
		}
		params[seg] = pathSegs[i]
	}
	return params
}
