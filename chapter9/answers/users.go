package answers

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// User represents a user.
type User struct {
	Key       *datastore.Key `datastore:"-"`
	Name      string
	AvatarURL string
}

// UserFromAppengineUser gets the User entity from the AppEngine User,
// or creates one if it doesn't exist.
func UserFromAppengineUser(ctx context.Context, u *user.User) (*User, error) {
	userKey := datastore.NewKey(ctx, "User", u.ID, 0, nil)
	var appUser User
	err := datastore.Get(ctx, userKey, &appUser)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}
	if err == datastore.ErrNoSuchEntity {
		appUser = User{
			Name:      u.String(),
			AvatarURL: gravatarURL(u.Email),
		}
		_, err := datastore.Put(ctx, userKey, &appUser)
		if err != nil {
			return nil, err
		}
	}
	appUser.Key = userKey
	return &appUser, nil
}

func (u User) AsUserEmbedded() UserEmbedded {
	return UserEmbedded{
		Key:       u.Key,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}
}

// UserEmbedded represents a User inside another object.
type UserEmbedded struct {
	Key       *datastore.Key
	Name      string
	AvatarURL string
}

func gravatarURL(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil))
}
