package answers

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.Handle("/error", handleError())
	http.Handle("/ask", MustLogin(handleAsk()))
	http.Handle("/questions/", handleQuestions())
	http.Handle("/answers/", handleAnswers())
	http.Handle("/", http.NotFoundHandler())
}

// PageData holds data for the page templates.
type PageData struct {
	Title     string
	User      *user.User
	LoginURL  string
	LogoutURL string
	Data      interface{}
}

// NewPageData gets a PageData object for page templates.
func NewPageData(ctx context.Context, r *http.Request, title string) PageData {
	logoutURL, err := user.LogoutURL(ctx, "/")
	if err != nil {
		log.Errorf(ctx, "LogoutURL: %s", err)
	}
	loginURL, err := user.LoginURL(ctx, r.URL.Path)
	if err != nil {
		log.Errorf(ctx, "LoginURL: %s", err)
	}
	return PageData{
		Title: title,
		// user stuff
		User:      user.Current(ctx),
		LogoutURL: logoutURL,
		LoginURL:  loginURL,
	}
}

// templateFiles gets the full list of template files from a list
// of names.
func templateFiles(names ...string) []string {
	names = append([]string{"base"}, names...)
	for i, name := range names {
		names[i] = filepath.Join("templates", name+".tmpl.html")
	}
	return names
}

func handleError() http.Handler {
	tmpl, err := template.ParseFiles(templateFiles("error")...)
	if err != nil {
		return ErrorHandler(err, http.StatusInternalServerError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		pageData := NewPageData(ctx, r, "Oops")
		pageData.Data = struct {
			ErrorMessage string
		}{
			ErrorMessage: r.URL.Query().Get("msg"),
		}
		err := tmpl.Execute(w, pageData)
		if err != nil {
			log.Errorf(ctx, "template: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// ErrorHandler is an http.Handler that always reports a specific
// error.
func ErrorHandler(err error, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		log.Errorf(ctx, "%d: %s", code, err)
		w.WriteHeader(code)
		io.WriteString(w, err.Error())
	})
}

// MustLogin wraps handlers forcing the user to login before they can
// access the handler.
func MustLogin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		if user.Current(ctx) == nil {
			loginURL, err := user.LoginURL(ctx, r.URL.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		}
		h.ServeHTTP(w, r)
	})
}

func RedirectToErrorPage(ctx context.Context, w http.ResponseWriter, r *http.Request, err string) {
	u, _ := url.Parse("/error")
	q := u.Query()
	q.Set("msg", err)
	u.RawQuery = q.Encode()
	log.Infof(ctx, "RedirectToErrorPage: %s", u.String())
	http.Redirect(w, r, u.String(), http.StatusFound)
}
