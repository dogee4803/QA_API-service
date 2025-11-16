package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"qna-api/internal/service"

	"github.com/gorilla/mux"
)

type QuestionHandler struct {
	questionService *service.QuestionService
}

func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// POST /questions
func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req service.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	question, err := h.questionService.CreateQuestion(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

// GET /questions
func (h *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	questions, err := h.questionService.GetAllQuestions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении вопросов")
		return
	}

	if questions == nil {
		questions = []service.Question{}
	}

	json.NewEncoder(w).Encode(questions)
}

// GET /questions/{id}
func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID вопроса")
		return
	}

	question, err := h.questionService.GetQuestion(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "вопрос не найден" {
			writeError(w, http.StatusNotFound, "Вопрос не найден")
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка при получении вопроса")
		}
		return
	}

	json.NewEncoder(w).Encode(question)
}

// DELETE /questions/{id}
func (h *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID вопроса")
		return
	}

	err = h.questionService.DeleteQuestion(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "вопрос не найден" {
			writeError(w, http.StatusNotFound, "Вопрос не найден")
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка при удалении вопроса")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}