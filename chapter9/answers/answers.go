package answers

import (
	"errors"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Answer represents an answer.
type Answer struct {
	Key    *datastore.Key `datastore:"-"`
	User   UserEmbedded
	Answer string
	Score  int
	CTime  time.Time
}

func NewAnswer(ctx context.Context, questionKey *datastore.Key, user *User, answer string) *Answer {
	answerKey := datastore.NewIncompleteKey(ctx, "Answer", questionKey)
	return &Answer{
		Key:    answerKey,
		User:   user.AsUserEmbedded(),
		Answer: answer,
		CTime:  time.Now(),
	}
}

func (a Answer) Valid() error {
	if len(a.Answer) < 10 {
		return errors.New("answer too short")
	}
	return nil
}

func (a Answer) NiceCTime() string {
	return humanize.Time(a.CTime)
}

func (a Answer) PositiveScore() bool {
	return a.Score >= 0
}

func (a Answer) UpVoteURL() string {
	return fmt.Sprintf("/answers/%s/up", a.Key.Encode())
}

func (a Answer) DownVoteURL() string {
	return fmt.Sprintf("/answers/%s/down", a.Key.Encode())
}

func (a *Answer) Put(ctx context.Context) error {
	var err error
	a.Key, err = datastore.Put(ctx, a.Key, a)
	return err
}

func (a *Answer) Vote(ctx context.Context, user *User, vote int) error {
	voteKey := VoteKey(ctx, a.Key, user.Key)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		// get or create the vote
		var voteEntity Vote
		var delta int
		err := datastore.Get(ctx, voteKey, &voteEntity)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		if err == datastore.ErrNoSuchEntity {
			voteEntity = Vote{
				AnswerKey: a.Key,
				UserKey:   user.Key,
				Value:     vote,
			}
		} else {
			// changing existing vote
			delta = -voteEntity.Value
		}
		voteEntity.Value = vote
		_, err = datastore.Put(ctx, voteKey, &voteEntity)
		if err != nil {
			return err
		}
		// update the answer
		var answer Answer
		err = datastore.Get(ctx, a.Key, &answer)
		if err != nil {
			return err
		}
		delta += vote
		answer.Score += delta
		_, err = datastore.Put(ctx, a.Key, &answer)
		if err != nil {
			return err
		}
		return nil // success
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		return err
	}
	return nil
}

func GetAnswer(ctx context.Context, answerKey *datastore.Key) (*Answer, error) {
	var answer Answer
	err := datastore.Get(ctx, answerKey, &answer)
	if err != nil {
		return nil, err
	}
	answer.Key = answerKey
	return &answer, nil
}

func GetAnswers(ctx context.Context, questionKey *datastore.Key) ([]*Answer, error) {
	var answers []*Answer
	keys, err := datastore.NewQuery("Answer").
		Ancestor(questionKey).
		Order("-Score").
		Order("-CTime").
		Limit(10).
		GetAll(ctx, &answers)
	if err != nil {
		return nil, err
	}
	for i, answer := range answers {
		answer.Key = keys[i]
	}
	return answers, nil
}
