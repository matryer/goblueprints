package answers

import (
	"html/template"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func handleQuestions() http.Handler {
	tmpl, err := template.ParseFiles(templateFiles("question.view")...)
	if err != nil {
		return ErrorHandler(err, http.StatusInternalServerError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		switch len(segs) {
		case 2:
			switch r.Method {
			case "GET":
				handleQuestionView(w, r, tmpl, segs[1])
			case "POST":
				handleSubmitAnswer(w, r, segs[1])
			}
		default:
			http.NotFound(w, r)
		}
	})
}

func handleQuestionView(w http.ResponseWriter, r *http.Request, tmpl *template.Template, questionID string) {
	ctx := appengine.NewContext(r)
	questionKey, err := datastore.DecodeKey(questionID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	question, err := GetQuestion(ctx, questionKey)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answers, err := GetAnswers(ctx, question.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageData := NewPageData(ctx, r, question.Text)
	pageData.Data = struct {
		Question *Question
		Answers  []*Answer
	}{
		Question: question,
		Answers:  answers,
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Errorf(ctx, "template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSubmitAnswer(w http.ResponseWriter, r *http.Request, questionID string) {
	ctx := appengine.NewContext(r)
	questionKey, err := datastore.DecodeKey(questionID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	u, err := UserFromAppengineUser(ctx, user.Current(ctx))
	if err != nil {
		log.Errorf(ctx, "UserFromAppengineUser(%s): %s", user.Current(ctx), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer := NewAnswer(ctx, questionKey, u, r.FormValue("answer"))
	err = answer.Valid()
	if err != nil {
		RedirectToErrorPage(ctx, w, r, err.Error())
		return
	}
	err = answer.Put(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.URL.Path, http.StatusFound)
}
