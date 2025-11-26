package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NKV510/question-answer-api/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateQuestionRequest struct {
	Text string `json:"text" binding:"required,min=1"`
}

func (h *Handler) GetQuestions(c *gin.Context) {
	ctx := c.Request.Context()

	slog.InfoContext(ctx, "Getting all questions")

	questions, err := h.repo.GetQuestions(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch questions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	ctx := c.Request.Context()

	slog.InfoContext(ctx, "Creating new question")

	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.ErrorContext(ctx, "Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	question := models.Question{
		Text: req.Text,
	}

	err := h.repo.CreateQuestion(ctx, &question)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create question", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusCreated, question)
}

func (h *Handler) GetQuestion(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid question ID", "error", err, "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	slog.InfoContext(ctx, "Getting question with answers", "question_id", id)

	question, err := h.repo.GetQuestion(ctx, uint(id))
	if err != nil {
		slog.ErrorContext(ctx, "Question not found", "question_id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid question ID", "error", err, "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	slog.InfoContext(ctx, "Deleting question", "question_id", id)

	err = h.repo.DeleteQuestion(ctx, uint(id))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete question", "question_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	c.Status(http.StatusNoContent)
}
