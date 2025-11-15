package repository

import (
	"context"
	"qna-api/internal/model"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	result := r.db.WithContext(ctx).Create(question)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]model.Question, error) {
	var questions []model.Question
	
	result := r.db.WithContext(ctx).Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return questions, nil
}

func (r *questionRepository) GetByID(ctx context.Context, id uint) (*model.Question, error) {
	var question model.Question
	
	result := r.db.WithContext(ctx).
		Preload("Answers").
		First(&question, id)
		
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, result.Error
	}
	
	return &question, nil
}

func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Question{}, id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}