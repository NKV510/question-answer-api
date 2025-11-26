package repository

import (
	"context"
	"log/slog"

	"github.com/NKV510/question-answer-api/internal/models"
)

func (r *Repository) CreateQuestion(ctx context.Context, question *models.Question) error {
	result := r.db.WithContext(ctx).Create(question)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to create question", "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *Repository) GetQuestions(ctx context.Context) ([]models.Question, error) {
	var questions []models.Question
	result := r.db.WithContext(ctx).Find(&questions)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to get questions", "error", result.Error)
		return nil, result.Error
	}
	return questions, nil
}

func (r *Repository) GetQuestion(ctx context.Context, id uint) (*models.Question, error) {
	var question models.Question
	result := r.db.WithContext(ctx).Preload("Answers").First(&question, id)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to get question", "id", id, "error", result.Error)
		return nil, result.Error
	}
	return &question, nil
}

func (r *Repository) DeleteQuestion(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Question{}, id)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to delete question", "id", id, "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *Repository) QuestionExists(ctx context.Context, id uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Question{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to check question existence", "id", id, "error", result.Error)
		return false, result.Error
	}
	return count > 0, nil
}
