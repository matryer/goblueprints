package main

import (
	"fmt"
	"log"
	"strings"

	"net/http"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated
		w.Header()["Location"] = []string{"/login"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error
		panic(err.Error())
	} else {
		// success - call the next handler
		h.next.ServeHTTP(w, r)
	}

}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":

		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("Error when trying to GetBeginAuthURL for", provider, "-", err)
		}

		w.Header()["Location"] = []string{loginUrl}
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":

		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}

		// get the credentials
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("Error when trying to complete auth for", provider, "-", err)
		}

		// get the user
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}

		// save some data
		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
			"email":      user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.Write([]byte(fmt.Sprintf("Auth action %s not supported", action)))
		w.WriteHeader(http.StatusNotFound)
	}
}
