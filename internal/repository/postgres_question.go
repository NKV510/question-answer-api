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
