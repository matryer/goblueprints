package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// templ represents a single template
type templ struct {
	source string
	templ  *template.Template
}

// Handle is a http.HandleFunc that renders this template.
func (t *templ) Handle(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.source)))
	}
	t.templ.Execute(w, nil)
}

func main() {

	// root
	http.HandleFunc("/", (&templ{source: "chat.html"}).Handle)

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
