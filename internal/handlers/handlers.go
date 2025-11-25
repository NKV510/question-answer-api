package handlers

import "github.com/NKV510/question-answer-api/internal/repository"

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
