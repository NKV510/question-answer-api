package handlers

import (
	"context"

	"github.com/NKV510/question-answer-api/internal/models"
)

// type Handler struct {
// 	repo *repository.Repository
// }

//	func NewHandler(repo *repository.Repository) *Handler {
//		return &Handler{
//			repo: repo,
//		}
//	}
type Repository interface {
	CreateQuestion(ctx context.Context, question *models.Question) error
	GetQuestions(ctx context.Context) ([]models.Question, error)
	GetQuestion(ctx context.Context, id uint) (*models.Question, error)
	DeleteQuestion(ctx context.Context, id uint) error
	QuestionExists(ctx context.Context, id uint) (bool, error)
	CreateAnswer(ctx context.Context, answer *models.Answer) error
	GetAnswer(ctx context.Context, id uint) (*models.Answer, error)
	DeleteAnswer(ctx context.Context, id uint) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
