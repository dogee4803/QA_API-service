package repository

import (
	"context"
	"qna-api/internal/model"

	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(ctx context.Context, answer *model.Answer) error {
	result := r.db.WithContext(ctx).Create(answer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *answerRepository) GetByID(ctx context.Context, id uint) (*model.Answer, error) {
	var answer model.Answer
	
	result := r.db.WithContext(ctx).First(&answer, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, result.Error
	}
	
	return &answer, nil
}

func (r *answerRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Answer{}, id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

func (r *answerRepository) GetByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error) {
	var answers []model.Answer
	
	result := r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Find(&answers)
		
	if result.Error != nil {
		return nil, result.Error
	}
	
	return answers, nil
}