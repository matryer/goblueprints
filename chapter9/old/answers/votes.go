package answers

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Vote struct {
	AnswerKey *datastore.Key
	UserKey   *datastore.Key
	Value     int
}

func VoteKey(ctx context.Context, answerKey, userKey *datastore.Key) *datastore.Key {
	pair := fmt.Sprintf("%s-%s", answerKey.Encode(), userKey.Encode())
	return datastore.NewKey(ctx, "Vote", pair, 0, nil)
}

func PutVote(ctx context.Context, answer *Answer, user *User, vote int) error {
	voteKey := VoteKey(ctx, answer.Key, user.Key)
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
				AnswerKey: answer.Key,
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
		err = datastore.Get(ctx, answer.Key, &answer)
		if err != nil {
			return err
		}
		delta += vote
		answer.Score += delta
		_, err = datastore.Put(ctx, answer.Key, &answer)
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
