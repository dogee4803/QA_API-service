package service

import (
	"context"
	"fmt"
	"qna-api/internal/model"
	"qna-api/internal/repository"
)

type QuestionService struct {
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
}

func NewQuestionService(questionRepo repository.QuestionRepository, answerRepo repository.AnswerRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
	}
}

type CreateQuestionRequest struct {
	Text string `json:"text"`
}

func (s *QuestionService) CreateQuestion(ctx context.Context, req *CreateQuestionRequest) (*model.Question, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("текст вопроса не может быть пустым")
	}

	if len(req.Text) > 1000 {
		return nil, fmt.Errorf("текст вопроса не может превышать 1000 символов")
	}

	question := &model.Question{
		Text: req.Text,
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, fmt.Errorf("ошибка при создании вопроса: %w", err)
	}

	return question, nil
}

func (s *QuestionService) GetAllQuestions(ctx context.Context) ([]model.Question, error) {
	questions, err := s.questionRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении вопросов: %w", err)
	}

	return questions, nil
}

func (s *QuestionService) GetQuestion(ctx context.Context, id uint) (*model.Question, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID вопроса не может быть пустым")
	}

	question, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("вопрос не найден: %w", err)
	}

	return question, nil
}

func (s *QuestionService) DeleteQuestion(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("ID вопроса не может быть пустым")
	}

	_, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("вопрос не найден: %w", err)
	}

	if err := s.questionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ошибка при удалении вопроса: %w", err)
	}

	return nil
}