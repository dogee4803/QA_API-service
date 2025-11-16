package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"qna-api/internal/service"

	"github.com/gorilla/mux"
)

type AnswerHandler struct {
	answerService *service.AnswerService
}

func NewAnswerHandler(answerService *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{
		answerService: answerService,
	}
}

// POST /questions/{id}/answers
func (h *AnswerHandler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	questionID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID вопроса")
		return
	}

	var req service.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	req.QuestionID = uint(questionID)

	answer, err := h.answerService.CreateAnswer(r.Context(), &req)
	if err != nil {
		errorMsg := err.Error()
		switch {
		case errorMsg == "вопрос не найден":
			writeError(w, http.StatusNotFound, "Вопрос не найден")
		case errorMsg == "некорректный формат ID пользователя":
			writeError(w, http.StatusBadRequest, errorMsg)
		default:
			writeError(w, http.StatusBadRequest, errorMsg)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(answer)
}

// GET /answers/{id}
func (h *AnswerHandler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID ответа")
		return
	}

	answer, err := h.answerService.GetAnswer(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "ответ не найден" {
			writeError(w, http.StatusNotFound, "Ответ не найден")
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка при получении ответа")
		}
		return
	}

	json.NewEncoder(w).Encode(answer)
}

// DELETE /answers/{id}
func (h *AnswerHandler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID ответа")
		return
	}

	err = h.answerService.DeleteAnswer(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "ответ не найден" {
			writeError(w, http.StatusNotFound, "Ответ не найден")
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка при удалении ответа")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}