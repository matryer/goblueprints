package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongoStorage is a private structure for
// storing polls with a MongoDB database
type mongoStorage struct {
	db *mgo.Session
}

// NewMongoStorage returns a Storage interface
// backed by the private mongoStorage struct
func NewMongoStorage(db *mgo.Session) PollStorage {
	return &mongoStorage{db: db}
}

func (s *mongoStorage) Get(p *Path) ([]*Poll, error) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	var q *mgo.Query
	if p.HasID() {
		// get specific poll
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all polls
		q = c.Find(nil)
	}
	var result []*Poll
	if err := q.All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *mongoStorage) Create(p *Poll) (string, error) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

func (s *mongoStorage) Delete(p *Path) error {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")

	if !p.HasID() {
		return ErrCannotDeleteAll
	}

	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		return err
	}

	return nil
}
