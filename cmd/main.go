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

func main() {
	setupLogging()

	slog.Info("Starting question-answer API server")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("Config loaded successfully",
		"db_host", cfg.DBHost,
		"db_port", cfg.DBPort,
		"db_name", cfg.DBName,
	)

	if err := database.ConnectDataBase(cfg); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(database.GetDB())

	handler := handlers.NewHandler(repo)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())

	setupRoutes(router, handler)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

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

func setupLogging() {
	if os.Getenv("ENV") == "production" {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
	} else {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger := slog.With(
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"client_ip", c.ClientIP(),
		)

		if c.Writer.Status() >= 500 {
			logger.Error("HTTP request error")
		} else if c.Writer.Status() >= 400 {
			logger.Warn("HTTP request client error")
		} else {
			logger.Info("HTTP request")
		}
	}
}

func setupRoutes(router *gin.Engine, handler *handlers.Handler) {
	questions := router.Group("/questions")
	{
		questions.GET("/", handler.GetQuestions)
		questions.POST("/", handler.CreateQuestion)
		questions.GET("/:id", handler.GetQuestion)
		questions.DELETE("/:id", handler.DeleteQuestion)
	}

	answers := router.Group("/answers")
	{
		answers.GET("/:id", handler.GetAnswer)
		answers.DELETE("/:id", handler.DeleteAnswer)
	}

	router.POST("/questions/:id/answers", handler.CreateAnswer)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "question-answer-api",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})
}
