package service

import (
	"context"
	"fmt"
	"qna-api/internal/model"
	"qna-api/internal/repository"

	"github.com/google/uuid"
)

type AnswerService struct {
	answerRepo   repository.AnswerRepository
	questionRepo repository.QuestionRepository
}

func NewAnswerService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

type CreateAnswerRequest struct {
	QuestionID uint   `json:"question_id"`
	UserID     string `json:"user_id"`
	Text       string `json:"text"`
}

func (s *AnswerService) CreateAnswer(ctx context.Context, req *CreateAnswerRequest) (*model.Answer, error) {

	_, err := s.questionRepo.GetByID(ctx, req.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("вопрос с ID %d не найден: %w", req.QuestionID, err)
	}

	if req.UserID == "" {
		return nil, fmt.Errorf("ID пользователя не может быть пустым")
	}
	
	if _, err := uuid.Parse(req.UserID); err != nil {
		return nil, fmt.Errorf("некорректный формат ID пользователя: %w", err)
	}

	if req.Text == "" {
		return nil, fmt.Errorf("текст ответа не может быть пустым")
	}

	if len(req.Text) > 2000 {
		return nil, fmt.Errorf("текст ответа не может превышать 2000 символов")
	}

	answer := &model.Answer{
		QuestionID: req.QuestionID,
		UserID:     req.UserID,
		Text:       req.Text,
	}

	if err := s.answerRepo.Create(ctx, answer); err != nil {
		return nil, fmt.Errorf("ошибка при создании ответа: %w", err)
	}

	return answer, nil
}

func (s *AnswerService) GetAnswer(ctx context.Context, id uint) (*model.Answer, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID ответа не может быть пустым")
	}

	answer, err := s.answerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ответ не найден: %w", err)
	}

	return answer, nil
}

func (s *AnswerService) DeleteAnswer(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("ID ответа не может быть пустым")
	}

	_, err := s.answerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ответ не найден: %w", err)
	}

	if err := s.answerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ошибка при удалении ответа: %w", err)
	}

	return nil
}