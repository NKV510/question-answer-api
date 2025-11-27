package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NKV510/question-answer-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockRepository) GetQuestions(ctx context.Context) ([]models.Question, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockRepository) GetQuestion(ctx context.Context, id uint) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}

func (m *MockRepository) DeleteQuestion(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateAnswer(ctx context.Context, answer *models.Answer) error {
	args := m.Called(ctx, answer)
	return args.Error(0)
}

func (m *MockRepository) GetAnswer(ctx context.Context, id uint) (*models.Answer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Answer), args.Error(1)
}

func (m *MockRepository) DeleteAnswer(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) QuestionExists(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestCreateQuestion_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.POST("/questions", handler.CreateQuestion)

	// Mock expectations
	mockRepo.On("CreateQuestion", mock.Anything, mock.AnythingOfType("*models.Question")).
		Run(func(args mock.Arguments) {
			question := args.Get(1).(*models.Question)
			question.ID = 1 // Симулируем присвоение ID базой данных
		}).
		Return(nil)

	// Test
	questionData := map[string]string{"text": "Test question?"}
	jsonData, _ := json.Marshal(questionData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test question?", response.Text)

	mockRepo.AssertExpectations(t)
}

func TestCreateQuestion_InvalidInput(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.POST("/questions", handler.CreateQuestion)

	// Test - empty text
	questionData := map[string]string{"text": ""}
	jsonData, _ := json.Marshal(questionData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверяем что метод репозитория не вызывался
	mockRepo.AssertNotCalled(t, "CreateQuestion")
}

func TestGetQuestions_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.GET("/questions", handler.GetQuestions)

	// Mock expectations
	expectedQuestions := []models.Question{
		{ID: 1, Text: "Question 1?"},
		{ID: 2, Text: "Question 2?"},
	}
	mockRepo.On("GetQuestions", mock.Anything).Return(expectedQuestions, nil)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, uint(1), response[0].ID)
	assert.Equal(t, "Question 1?", response[0].Text)

	mockRepo.AssertExpectations(t)
}

func TestGetQuestions_Error(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.GET("/questions", handler.GetQuestions)

	// Mock expectations - возвращаем ошибку
	mockRepo.On("GetQuestions", mock.Anything).Return([]models.Question{}, assert.AnError)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Failed to fetch questions")

	mockRepo.AssertExpectations(t)
}

func TestGetQuestion_Success(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.GET("/questions/:id", handler.GetQuestion)

	expectedQuestion := &models.Question{
		ID:   1,
		Text: "Test question?",
		Answers: []models.Answer{
			{ID: 1, QuestionID: 1, UserID: "user1", Text: "Answer 1"},
		},
	}
	mockRepo.On("GetQuestion", mock.Anything, uint(1)).Return(expectedQuestion, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test question?", response.Text)
	assert.Len(t, response.Answers, 1)

	mockRepo.AssertExpectations(t)
}

func TestGetQuestion_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.GET("/questions/:id", handler.GetQuestion)

	mockRepo.On("GetQuestion", mock.Anything, uint(999)).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Question not found", response["error"])

	mockRepo.AssertExpectations(t)
}

func TestGetQuestion_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(MockRepository)
	handler := NewHandler(mockRepo)

	router.GET("/questions/:id", handler.GetQuestion)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid question ID")

	mockRepo.AssertNotCalled(t, "GetQuestion")
}
