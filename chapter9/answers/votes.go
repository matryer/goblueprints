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
