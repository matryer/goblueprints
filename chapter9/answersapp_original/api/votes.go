package api

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Vote struct {
	Key      *datastore.Key `json:"id" datastore:"-"`
	MTime    time.Time      `json:"last_modified"`
	Question QuestionCard   `json:"question"`
	Answer   AnswerCard     `json:"answer"`
	User     UserCard       `json:"user"`
	Score    int            `json:"score"`
}

func (v *Vote) Put(ctx context.Context) error {
	var err error
	v.Key, err = datastore.Put(ctx, v.Key, v)
	if err != nil {
		return err
	}
	return nil
}

func CastVote(ctx context.Context, answerKey *datastore.Key, score int) (*Vote, error) {
	user, err := UserFromAEUser(ctx)
	if err != nil {
		return nil, err
	}
	question, err := GetQuestion(ctx, answerKey.Parent())
	if err != nil {
		return nil, err
	}
	// key is made up of answer key, plus user key
	voteKeyStr := fmt.Sprintf("%s:%s", answerKey.Encode(), user.Key.Encode())
	voteKey := datastore.NewKey(ctx, "Vote", voteKeyStr, 0, nil)
	var vote Vote
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		answer, err := GetAnswer(ctx, answerKey)
		if err != nil {
			return err
		}
		var delta int // delta describes the change to answer score
		err = datastore.Get(ctx, voteKey, &vote)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		if err == datastore.ErrNoSuchEntity {
			vote = Vote{
				Key:      voteKey,
				User:     user.Card(),
				Answer:   answer.Card(),
				Question: question.Card(),
				Score:    score,
			}
		} else {
			// they have already voted - so we will be changing
			// this vote
			delta = vote.Score * -1
		}
		delta += score
		answer.Score += delta
		err = answer.Put(ctx)
		if err != nil {
			return err
		}
		vote.Key = voteKey
		vote.Score = score
		vote.MTime = time.Now()
		err = vote.Put(ctx)
		if err != nil {
			return err
		}
		return nil

	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func validScore(score int) error {
	if score != -1 && score != 1 {
		return errors.New("invalid score")
	}
	return nil
}
