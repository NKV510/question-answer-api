package repository

import (
	"context"
	"log/slog"

	"github.com/NKV510/question-answer-api/internal/models"
)

func (r *Repository) CreateAnswer(ctx context.Context, answer *models.Answer) error {
	result := r.db.WithContext(ctx).Create(answer)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to create answer", "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *Repository) GetAnswer(ctx context.Context, id uint) (*models.Answer, error) {
	var answer models.Answer
	result := r.db.WithContext(ctx).First(&answer, id)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to get answer", "id", id, "error", result.Error)
		return nil, result.Error
	}
	return &answer, nil
}
