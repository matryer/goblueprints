package answers

import (
	"html/template"
	"net/http"
	"strings"

	"google.golang.org/appengine/user"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func handleAnswers() http.Handler {
	tmpl, err := template.ParseFiles(templateFiles("answer")...)
	if err != nil {
		return ErrorHandler(err, http.StatusInternalServerError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		switch len(segs) {
		case 2:
			handleAnswerView(w, r, tmpl, segs[1])
		case 3:
			handleAnswerVote(w, r, segs[1], segs[2])
		default:
			http.NotFound(w, r)
		}
	})
}

func handleAnswerView(w http.ResponseWriter, r *http.Request, tmpl *template.Template, answerID string) {
	ctx := appengine.NewContext(r)
	answerKey, err := datastore.DecodeKey(answerID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	answer, err := GetAnswer(ctx, answerKey)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	question, err := GetQuestion(ctx, answer.Key.Parent())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageData := NewPageData(ctx, r, question.Text+" via "+answer.User.Name)
	pageData.Data = struct {
		Answer   *Answer
		Question *Question
	}{
		Answer:   answer,
		Question: question,
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Errorf(ctx, "template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAnswerVote(w http.ResponseWriter, r *http.Request, answerID, upOrDown string) {
	ctx := appengine.NewContext(r)
	var vote int
	switch upOrDown {
	case "up":
		vote = 1
	case "down":
		vote = -1
	default:
		http.NotFound(w, r)
		return
	}
	usr := user.Current(ctx)
	if usr == nil {
		loginURL, err := user.LoginURL(ctx, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		return
	}
	answerKey, err := datastore.DecodeKey(answerID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	answer, err := GetAnswer(ctx, answerKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := UserFromAppengineUser(ctx, usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = answer.Vote(ctx, user, vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
