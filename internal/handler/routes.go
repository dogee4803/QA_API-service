package handler

import (
	"net/http"
	"qna-api/internal/service"

	"github.com/gorilla/mux"
)

func SetupRoutes(questionService *service.QuestionService, answerService *service.AnswerService) *mux.Router {
	router := mux.NewRouter()

	questionHandler := NewQuestionHandler(questionService)
	answerHandler := NewAnswerHandler(answerService)

	router.HandleFunc("/questions", questionHandler.GetQuestions).Methods("GET")
	router.HandleFunc("/questions", questionHandler.CreateQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", questionHandler.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", questionHandler.DeleteQuestion).Methods("DELETE")

	router.HandleFunc("/questions/{id}/answers", answerHandler.CreateAnswer).Methods("POST")
	router.HandleFunc("/answers/{id}", answerHandler.GetAnswer).Methods("GET")
	router.HandleFunc("/answers/{id}", answerHandler.DeleteAnswer).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"service":   "qna-api",
			"version":   "1.0.0",
		})
	}).Methods("GET")

	return router
}