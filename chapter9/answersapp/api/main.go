package api

import "net/http"

func init() {
	http.HandleFunc("/api/questions/", handleQuestions)
	http.HandleFunc("/api/answers/", handleAnswers)
	http.HandleFunc("/api/votes/", handleVotes)
}
