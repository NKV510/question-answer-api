package handlers

import (
	"log/slog"
	"net/http"

	"github.com/NKV510/question-answer-api/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetQuestionHandler(c *gin.Context) {

}

func (h *Handler) PostQuestionHandler(c *gin.Context) {

	ctx := c.Request.Context()

	slog.InfoContext(ctx, "Creating new question")

	var req struct {
		Text string `json:"text" binding:"required,min=1"`
	}
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
