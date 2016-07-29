package answers

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func handleAsk() http.Handler {
	tmpl, err := template.ParseFiles(templateFiles("ask")...)
	if err != nil {
		return ErrorHandler(err, http.StatusInternalServerError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleAskGet(w, r, tmpl)
		case "POST":
			handleAskPost(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func handleAskGet(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	ctx := appengine.NewContext(r)
	pageData := NewPageData(ctx, r, "Ask a question")
	err := tmpl.Execute(w, pageData)
	if err != nil {
		log.Errorf(ctx, "template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAskPost(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u, err := UserFromAppengineUser(ctx, user.Current(ctx))
	if err != nil {
		log.Errorf(ctx, "UserFromAppengineUser(%s): %s", user.Current(ctx), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	question := NewQuestion(ctx, u, r.FormValue("question"))
	err = question.Valid()
	if err != nil {
		RedirectToErrorPage(ctx, w, r, err.Error())
		return
	}
	err = question.Put(ctx)
	if err != nil {
		RedirectToErrorPage(ctx, w, r, fmt.Sprintf("couldn't save question: %s", err))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("questions/%s", question.Key.Encode()), http.StatusFound)
}
