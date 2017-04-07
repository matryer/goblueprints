package main

import (
	"io/ioutil"
	"os"
	"path"

	gomniauthtest "github.com/stretchr/gomniauth/test"

	"testing"
)

func TestAuthAvatar(t *testing.T) {

	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testURL := "http://url-to-gravatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testURL, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}
func TestGravatarAvatar(t *testing.T) {

	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}

	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}

}

func TestFileSystemAvatar(t *testing.T) {

	// make a test avatar file
	filename := path.Join("avatars", "abc.jpg")
	if err := os.MkdirAll("avatars", 0777); err != nil {
		t.Errorf("couldn't make avatar dir: %s", err)
	}
	if err := ioutil.WriteFile(filename, []byte{}, 0777); err != nil {
		t.Errorf("couldn't make avatar: %s", err)
	}
	defer os.Remove(filename)

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}

	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Errorf("FileSystemAvatar.GetAvatarURL should not return an error: %s", err)
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}

}
