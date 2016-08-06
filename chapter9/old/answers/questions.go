package answers

import (
	"errors"
	"time"

	humanize "github.com/dustin/go-humanize"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Question represents a single question.
type Question struct {
	Key   *datastore.Key `datastore:"-"`
	User  UserEmbedded
	Text  string
	Tags  []string
	CTime time.Time
}

func NewQuestion(ctx context.Context, user *User, text string) *Question {
	questionKey := datastore.NewIncompleteKey(ctx, "Question", nil)
	return &Question{
		Key:   questionKey,
		User:  user.AsUserEmbedded(),
		Text:  text,
		CTime: time.Now(),
	}
}

func (q Question) NiceCTime() string {
	return humanize.Time(q.CTime)
}

func (q Question) Valid() error {
	if len(q.Text) < 10 {
		return errors.New("question is too short")
	}
	return nil
}

func (q *Question) Put(ctx context.Context) error {
	var err error
	q.Key, err = datastore.Put(ctx, q.Key, q)
	return err
}

func GetQuestion(ctx context.Context, questionKey *datastore.Key) (*Question, error) {
	var q Question
	err := datastore.Get(ctx, questionKey, &q)
	if err != nil {
		return nil, err
	}
	q.Key = questionKey
	return &q, nil
}
