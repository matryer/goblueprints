package main

import "testing"

func TestAuthAvatar(t *testing.T) {

	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	// set a value
	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}

	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct URL")
		}
	}

}

func TestGravitarAvatar(t *testing.T) {

	var gravatarAvitar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"email": "MyEmailAddress@example.com"}

	url, err := gravatarAvitar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvitar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
	}

}
