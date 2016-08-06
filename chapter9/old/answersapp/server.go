package answers

import "net/http"

type QuestionStore interface {
}

type Server struct {
	questions QuestionStore
}

func NewServer(questions QuestionStore) *Server {
	return &Server{
		questions: questions,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
