package api

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Question struct {
	Key          *datastore.Key `json:"id" datastore:"-"`
	CTime        time.Time      `json:"created"`
	Question     string         `json:"question"`
	User         UserCard       `json:"user"`
	AnswersCount int            `json:"answers_count"`
}

func (q Question) OK() error {
	if len(q.Question) < 10 {
		return errors.New("question is too short")
	}
	return nil
}

func (q *Question) Put(ctx context.Context) error {
	log.Debugf(ctx, "Saving question")
	if q.Key == nil {
		q.Key = datastore.NewIncompleteKey(ctx, "Question", nil)
	}
	user, err := UserFromAEUser(ctx)
	if err != nil {
		return err
	}
	log.Debugf(ctx, "UserFromAEUser: %v", user)
	q.User = user.Card()
	q.CTime = time.Now()
	q.Key, err = datastore.Put(ctx, q.Key, q)
	if err != nil {
		return err
	}
	return nil
}

func (q *Question) NewAnswer(ctx context.Context, answer string) *Answer {
	return &Answer{
		Key:    datastore.NewIncompleteKey(ctx, "Answer", q.Key),
		Answer: answer,
	}
}

func (q *Question) Card() QuestionCard {
	return QuestionCard{
		Key:      q.Key,
		Question: q.Question,
		User:     q.User,
	}
}

// GetQuestion gets a Question by key.
func GetQuestion(ctx context.Context, key *datastore.Key) (*Question, error) {
	var q Question
	err := datastore.Get(ctx, key, &q)
	if err != nil {
		return nil, err
	}
	q.Key = key
	return &q, nil
}

func TopQuestions(ctx context.Context) ([]*Question, error) {
	var questions []*Question
	questionKeys, err := datastore.NewQuery("Question").
		Order("-Score").
		Order("-CTime").
		Limit(100).
		GetAll(ctx, &questions)
	if err != nil {
		return nil, err
	}
	for i := range questions {
		questions[i].Key = questionKeys[i]
	}
	return questions, nil
}

type QuestionCard struct {
	Key      *datastore.Key `json:"id" datastore:"-"`
	Question string         `json:"question"`
	User     UserCard       `json:"user"`
}
