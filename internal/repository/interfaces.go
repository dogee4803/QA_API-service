package repository

import (
	"context"
	"qna-api/internal/model"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	
	GetAll(ctx context.Context) ([]model.Question, error)
	
	GetByID(ctx context.Context, id uint) (*model.Question, error)
	
	Delete(ctx context.Context, id uint) error
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) error
	
	GetByID(ctx context.Context, id uint) (*model.Answer, error)
	
	Delete(ctx context.Context, id uint) error
	
	GetByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error)
}