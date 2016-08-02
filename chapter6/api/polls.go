package main

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrCannotDeleteAll is used to identify a client-facing
	// error when deleting a poll from storage
	ErrCannotDeleteAll = errors.New("cannot delete all polls")
)

// Poll represents a poll. ID should not be a BSON object
type Poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title" bson:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey  string         `json:"apikey"`
}

// PollStorage represents an abstraction to a storage
// backend for Polls
type PollStorage interface {
	Get(p *Path) ([]*Poll, error)
	Create(p *Poll) (string, error)
	Delete(p *Path) error
}
