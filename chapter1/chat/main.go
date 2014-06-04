package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	source string
	templ  *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.source)))
	}
	t.templ.Execute(w, nil)
}

func main() {

	// root
	http.Handle("/", &templateHandler{source: "chat.html"})

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
