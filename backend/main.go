package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"library-management-system/internal/api"
	"library-management-system/internal/config"
	"library-management-system/internal/middleware"
	"library-management-system/internal/repository"
	"library-management-system/internal/service"
	"library-management-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Library Management System API")

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := repository.NewDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowingRepo := repository.NewBorrowingRepository(db)
	authLogRepo := repository.NewAuthLogRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, authLogRepo, tokenRepo, cfg.Auth)
	userService := service.NewUserService(userRepo)
	bookService := service.NewBookService(bookRepo)
	borrowingService := service.NewBorrowingService(borrowingRepo, bookRepo, userRepo)

	// Initialize router
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.RateLimiter(cfg.RateLimit))
	router.Use(middleware.SecurityHeaders())

	// Initialize API handlers
	api.RegisterRoutes(router, authService, userService, bookService, borrowingService, cfg)

	// Setup HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server is starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeoutSeconds)*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}
