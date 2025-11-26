package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NKV510/question-answer-api/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateAnswerRequest struct {
	UserID string `json:"user_id" binding:"required,min=1"`
	Text   string `json:"text" binding:"required,min=1"`
}

func (h *Handler) CreateAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid question ID", "error", err, "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	slog.InfoContext(ctx, "Creating answer for question", "question_id", questionID)

	exists, err := h.repo.QuestionExists(ctx, uint(questionID))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check question existence", "question_id", questionID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		slog.WarnContext(ctx, "Question not found", "question_id", questionID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.ErrorContext(ctx, "Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	answer := models.Answer{
		QuestionID: uint(questionID),
		UserID:     req.UserID,
		Text:       req.Text,
	}

	err = h.repo.CreateAnswer(ctx, &answer)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create answer", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, answer)
}

func (h *Handler) GetAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid answer ID", "error", err, "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	slog.InfoContext(ctx, "Getting answer", "answer_id", id)

	answer, err := h.repo.GetAnswer(ctx, uint(id))
	if err != nil {
		slog.ErrorContext(ctx, "Answer not found", "answer_id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, answer)
}

func (h *Handler) DeleteAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid answer ID", "error", err, "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	slog.InfoContext(ctx, "Deleting answer", "answer_id", id)

	err = h.repo.DeleteAnswer(ctx, uint(id))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete answer", "answer_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer"})
		return
	}

	c.Status(http.StatusNoContent)
}
