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
	MTime    time.Time      `json:"last_modified" datastore:",noindex"`
	Question QuestionCard   `json:"question" datastore:",noindex"`
	Answer   AnswerCard     `json:"answer" datastore:",noindex"`
	User     UserCard       `json:"user" datastore:",noindex"`
	Score    int            `json:"score" datastore:",noindex"`
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
	question, err := GetQuestion(ctx, answerKey.Parent())
	if err != nil {
		return nil, err
	}
	user, err := UserFromAEUser(ctx)
	if err != nil {
		return nil, err
	}
	var vote Vote
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var err error
		vote, err = castVoteInTransaction(ctx, answerKey, question, user, score)
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

func castVoteInTransaction(ctx context.Context,
	answerKey *datastore.Key,
	question *Question,
	user *User,
	score int) (Vote, error) {
	var vote Vote
	answer, err := GetAnswer(ctx, answerKey)
	if err != nil {
		return vote, err
	}
	voteKeyStr := fmt.Sprintf("%s:%s", answerKey.Encode(), user.Key.Encode())
	voteKey := datastore.NewKey(ctx, "Vote", voteKeyStr, 0, nil)
	var delta int // delta describes the change to answer score
	err = datastore.Get(ctx, voteKey, &vote)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return vote, err
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
		return vote, err
	}
	vote.Key = voteKey
	vote.Score = score
	vote.MTime = time.Now()
	err = vote.Put(ctx)
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func validScore(score int) error {
	if score != -1 && score != 1 {
		return errors.New("invalid score")
	}
	return nil
}
