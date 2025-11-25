package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NKV510/question-answer-api/internal/config"
	"github.com/NKV510/question-answer-api/internal/database"
	"github.com/NKV510/question-answer-api/internal/handlers"
	"github.com/NKV510/question-answer-api/internal/repository"
	"github.com/gin-gonic/gin"
)

const (
	envlocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.LoadConfig()

	log := setupLoger(cfg.ENV)

	log.Info("Start server", slog.String("env", cfg.ENV))

	if err != nil {
		log.Error("Can not config connection")
		os.Exit(1)
	}

	// Инициализация репозитория
	repo := repository.NewRepository(database.GetDB())

	// Инициализация обработчиков
	handler := handlers.NewHandler(repo)

	// Настройка Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())

	// Routes
	setupRoutes(router, handler)

	// Настройка сервера
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		slog.Info("Server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}

func setupLoger(ENV string) *slog.Logger {
	var log *slog.Logger
	switch ENV {
	case envlocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupRoutes(router *gin.Engine, handler *handlers.Handler) {
	// Questions endpoints
	questions := router.Group("/questions")
	{
		//questions.GET("/", handler.GetQuestions)
		questions.POST("/", handler.PostQuestionHandler)
		//questions.GET("/:id", handler.GetQuestion)
		//questions.DELETE("/:id", handler.DeleteQuestion)
	}

	// Answers endpoints
	// answers := router.Group("/answers")
	// {
	// 	answers.GET("/:id", handler.GetAnswer)
	// 	answers.DELETE("/:id", handler.DeleteAnswer)
	// }

	// // Ответы к конкретному вопросу
	// router.POST("/questions/:id/answers", handler.CreateAnswer)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "question-answer-api",
		})
	})

	// Обработка несуществующих маршрутов
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})
}
