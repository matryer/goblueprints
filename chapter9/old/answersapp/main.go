package answers

import (
	"net/http"

	"github.com/matryer/goblueprints/chapter9/answersapp/gae"
)

func init() {
	store := gae.NewQuestionStore()
	server := NewServer(store)
	http.Handle("/", server)
}
`