package tinypins

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "tasks/clean" {
		handleCleanup(w, r)
		return
	}
	segs := strings.Split(path, "/")
	switch len(segs) {
	case 1: // /
		handleLanding(w, r)
	case 2: // /pins/ID
		handlePins(w, r)
	case 3: // /pins/ID/updates
		handlePinsUpdates(w, r)
	default: // everything else
		http.NotFound(w, r)
	}
}

func handleLanding(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewV4()
	idStr := strings.Replace(id.String(), "-", "", -1)
	http.Redirect(w, r, fmt.Sprintf("/pins/%s", idStr), http.StatusFound)
}

func handlePins(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	me := user.Current(ctx)
	if me == nil {
		loginURL, err := user.LoginURL(ctx, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		return
	}
	http.ServeFile(w, r, "pages/pins.html")
}

func handlePinsUpdates(w http.ResponseWriter, r *http.Request) {
	switch strings.ToLower(r.Method) {
	case "get":
		handlePinsGetUpdate(w, r)
	case "post":
		handlePinsPostUpdate(w, r)
	}
}

// Pin represents a user at a location.
type Pin struct {
	UserID       string             `json:"user_id"`
	Name         string             `json:"name"`
	ImageURL     string             `json:"image_url"`
	IsMe         bool               `json:"me" datastore:"-"`
	Accuracy     float64            `json:"accuracy"`
	Location     appengine.GeoPoint `json:"location"`
	LastModified time.Time          `json:"last_modified"`
	SecondsAgo   float64            `json:"seconds_ago" datastore:"-"`
}

func handlePinsGetUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	me := user.Current(ctx)
	if me == nil {
		http.Error(w, "must be logged in", http.StatusUnauthorized)
		return
	}
	id, err := idFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ancestorKey := datastore.NewKey(ctx, "Groups", id, 0, nil)
	var pins []*Pin
	_, err = datastore.NewQuery("Pins").Ancestor(ancestorKey).GetAll(ctx, &pins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, pin := range pins {
		pin.IsMe = pin.UserID == me.ID
		pin.SecondsAgo = time.Now().Sub(pin.LastModified).Seconds()
	}
	err = json.NewEncoder(w).Encode(pins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePinsPostUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	me := user.Current(ctx)
	if me == nil {
		http.Error(w, "must be logged in", http.StatusUnauthorized)
		return
	}
	id, err := idFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var pin Pin
	err = json.NewDecoder(r.Body).Decode(&pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !pin.Location.Valid() {
		http.Error(w, "bad location", http.StatusBadRequest)
		return
	}
	pin.IsMe = true
	pin.UserID = me.ID
	pin.ImageURL = gravatarURL(me.Email)
	pin.Name = me.String()
	pin.LastModified = time.Now()
	ancestorKey := datastore.NewKey(ctx, "Groups", id, 0, nil)
	pinKey := datastore.NewKey(ctx, "Pins", me.ID, 0, ancestorKey)
	_, err = datastore.Put(ctx, pinKey, &pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCleanup(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Appengine-Cron") != "true" {
		http.NotFound(w, r)
		return
	}
	ctx := appengine.NewContext(r)
	keys, err := datastore.NewQuery("Pins").
		Filter("LastModified <", time.Now().Add(-1*time.Hour)).
		KeysOnly().GetAll(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "cleaning up %d item(s)", len(keys))
	err = datastore.DeleteMulti(ctx, keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Deleted int `json:"deleted"`
	}{Deleted: len(keys)})
}

func idFromRequest(r *http.Request) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	segs := strings.Split(path, "/")
	if len(segs) < 2 {
		return "", errors.New("missing ID")
	}
	if len(segs[1]) < 10 {
		return "", errors.New("invalid ID")
	}
	return segs[1], nil
}

func gravatarURL(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil))
}
