package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/guilherme/pet-services-api/internal/domain/auth"
	"github.com/guilherme/pet-services-api/internal/domain/user"
	infradatabase "github.com/guilherme/pet-services-api/internal/infra/database"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
	httpapi "github.com/guilherme/pet-services-api/internal/transport/http"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.Default()

	// ─────────────────────────────────────────────────────────────────────────
	// 1. Configuration
	// ─────────────────────────────────────────────────────────────────────────

	cfg := infradatabase.DefaultConfig()
	cfg.Logger = logger

	// ─────────────────────────────────────────────────────────────────────────
	// 2. Database Connection
	// ─────────────────────────────────────────────────────────────────────────

	db, sqlDB, err := infradatabase.Open(ctx, cfg)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer infradatabase.Close(logger, sqlDB)
	logger.Info("database connected")

	// ─────────────────────────────────────────────────────────────────────────
	// 3. Infrastructure Providers (placeholder implementations)
	// ─────────────────────────────────────────────────────────────────────────

	// TODO: Replace with real implementations (bcrypt, JWT, email service, etc.)
	tokenService := &stubTokenService{}
	passwordHasher := &stubPasswordHasher{}
	emailService := &stubEmailService{}

	// ─────────────────────────────────────────────────────────────────────────
	// 4. Factory: Use Cases + Repositories
	// ─────────────────────────────────────────────────────────────────────────

	useCases := factory.NewUseCases(factory.Config{
		DB:                   db,
		TokenService:         tokenService,
		PasswordHasher:       passwordHasher,
		EmailService:         emailService,
		PasswordResetBaseURL: os.Getenv("PASSWORD_RESET_BASE_URL"),
		EmailVerifyBaseURL:   os.Getenv("EMAIL_VERIFY_BASE_URL"),
		Logger:               logger,
	})

	// ─────────────────────────────────────────────────────────────────────────
	// 5. HTTP Router + Middlewares
	// ─────────────────────────────────────────────────────────────────────────

	router := gin.Default()

	// Global middlewares
	router.Use(requestIDMiddleware())
	router.Use(structuredLoggingMiddleware(logger))
	router.Use(corsMiddleware())

	// Health endpoints (public)
	router.GET("/health", healthHandler(db))
	router.GET("/ready", readinessHandler(db))

	// Register v1 routes
	httpapi.NewRouter(useCases, tokenService)
	routerWithRoutes := httpapi.NewRouter(useCases, tokenService)

	// Copy routes from the new router to our router with middlewares
	for _, route := range routerWithRoutes.Routes() {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 6. HTTP Server
	// ─────────────────────────────────────────────────────────────────────────

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 7. Graceful Shutdown Setup
	// ─────────────────────────────────────────────────────────────────────────

	errChan := make(chan error, 1)

	go func() {
		logger.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("server error", "error", err)
		os.Exit(1)
	case sig := <-sigChan:
		logger.Info("signal received", "signal", sig.String())

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown error", "error", err)
			os.Exit(1)
		}

		logger.Info("server shutdown complete")
	}
}

// ─────────────────────────────────────────────────────────────────────────
// Middlewares
// ─────────────────────────────────────────────────────────────────────────

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func structuredLoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, _ := c.Get("request_id")

		c.Next()

		duration := time.Since(start)
		logger.Info("http request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := os.Getenv("CORS_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ─────────────────────────────────────────────────────────────────────────
// Handlers
// ─────────────────────────────────────────────────────────────────────────

func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func readinessHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"reason": "database pool error",
			})
			return
		}

		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"reason": "database unreachable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────
// Stub Implementations (TODO: Replace with real implementations)
// ─────────────────────────────────────────────────────────────────────────

type stubTokenService struct{}

func (s *stubTokenService) GenerateTokens(userID uuid.UUID, userType user.UserType) (auth.TokenPair, error) {
	return auth.TokenPair{
		AccessToken:      fmt.Sprintf("access_%s_%d", userID.String()[:8], time.Now().Unix()),
		AccessExpiresAt:  time.Now().Add(15 * time.Minute),
		RefreshToken:     fmt.Sprintf("refresh_%s_%d", userID.String()[:8], time.Now().Unix()),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		RefreshID:        uuid.New(),
	}, nil
}

func (s *stubTokenService) ParseRefreshToken(token string) (auth.RefreshClaims, error) {
	return auth.RefreshClaims{
		TokenID:   uuid.New(),
		UserID:    uuid.New(),
		UserType:  user.UserTypeOwner,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

type stubPasswordHasher struct{}

func (s *stubPasswordHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (s *stubPasswordHasher) Compare(hash, password string) error {
	if hash != "hashed_"+password {
		return fmt.Errorf("password mismatch")
	}
	return nil
}

type stubEmailService struct{}

func (s *stubEmailService) SendPasswordResetEmail(to, resetLink string) error {
	fmt.Printf("📧 Password reset email sent to %s: %s\n", to, resetLink)
	return nil
}

func (s *stubEmailService) SendEmailVerification(to, verificationLink string) error {
	fmt.Printf("📧 Email verification sent to %s: %s\n", to, verificationLink)
	return nil
}
