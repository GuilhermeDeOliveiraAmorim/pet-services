package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pet-services-api/internal/application/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	infrabcrypt "pet-services-api/internal/infra/bcrypt"
	infradatabase "pet-services-api/internal/infra/database"
	infraemail "pet-services-api/internal/infra/email"
	"pet-services-api/internal/infra/factory"
	infrajwt "pet-services-api/internal/infra/jwt"
	httpapi "pet-services-api/internal/transport/http"
)

// @title Pet Services API
// @version 1.0
// @description This is an API for managing pet services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.petland.com.br
// @contact.email contato@petland.com.br

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1/
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.NewSlogLogger()

	// ─────────────────────────────────────────────────────────────────────────
	// 1. Configuration
	// ─────────────────────────────────────────────────────────────────────────

	cfg := infradatabase.DefaultConfig()
	cfg.Logger = logger.Logger() // *slog.Logger para o banco

	// ─────────────────────────────────────────────────────────────────────────
	// 2. Database Connection
	// ─────────────────────────────────────────────────────────────────────────

	db, sqlDB, err := infradatabase.Open(ctx, cfg)
	if err != nil {
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    500,
			Message: "failed to connect database",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.ERROR,
			Error:   err,
		})
		os.Exit(1)
	}
	defer infradatabase.Close(logger.Logger(), sqlDB)
	logger.Log(logging.Logger{
		Context: ctx,
		Code:    200,
		Message: "database connected",
		From:    "main",
		Layer:   logging.LoggerLayers.SERVER,
		TypeLog: logging.LoggerTypes.INFO,
	})

	// ─────────────────────────────────────────────────────────────────────────
	// 3. Infrastructure Providers (real implementations)
	// ─────────────────────────────────────────────────────────────────────────

	// JWT Token Service
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    500,
			Message: "JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.ERROR,
		})
		os.Exit(1)
	}

	accessDuration := parseDuration(os.Getenv("JWT_ACCESS_DURATION"), 15*time.Minute)
	refreshDuration := parseDuration(os.Getenv("JWT_REFRESH_DURATION"), 7*24*time.Hour)

	tokenService := infrajwt.NewTokenService(infrajwt.Config{
		AccessSecret:    accessSecret,
		RefreshSecret:   refreshSecret,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	})
	logger.Log(logging.Logger{
		Context: ctx,
		Code:    200,
		Message: "jwt token service initialized",
		From:    "main",
		Layer:   logging.LoggerLayers.SERVER,
		TypeLog: logging.LoggerTypes.INFO,
	})

	// Password Hasher (Bcrypt)
	passwordHasher := infrabcrypt.NewPasswordHasher()
	logger.Log(logging.Logger{
		Context: ctx,
		Code:    200,
		Message: "bcrypt password hasher initialized",
		From:    "main",
		Layer:   logging.LoggerLayers.SERVER,
		TypeLog: logging.LoggerTypes.INFO,
	})

	// Email Service (SMTP or Stub)
	emailHost := os.Getenv("SMTP_HOST")
	emailPort := os.Getenv("SMTP_PORT")
	emailUser := os.Getenv("SMTP_USER")
	emailPass := os.Getenv("SMTP_PASS")
	emailFrom := os.Getenv("SMTP_FROM")

	var emailService infraemail.EmailServiceInterface
	if strings.TrimSpace(emailHost) == "" || strings.TrimSpace(emailHost) == "localhost" {
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    200,
			Message: "using stub email service (no SMTP configured)",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.INFO,
		})
		emailService = infraemail.NewStubEmailService(logger.Logger())
	} else {
		port := 587
		if emailPort != "" {
			if p, err := strconv.Atoi(emailPort); err == nil {
				port = p
			}
		}
		emailService = infraemail.NewSMTPService(infraemail.Config{
			Host:     emailHost,
			Port:     port,
			User:     emailUser,
			Password: emailPass,
			FromAddr: emailFrom,
			Logger:   logger.Logger(),
		})
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    200,
			Message: "smtp email service initialized",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.INFO,
		})
	}

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
	router.Use(structuredLoggingMiddleware(logger.Logger()))
	router.Use(corsMiddleware())

	// Health endpoints (public)
	router.GET("/health", healthHandler(db))
	router.GET("/ready", readinessHandler(db))

	// Register v1 routes
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
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    200,
			Message: "server starting",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.INFO,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    500,
			Message: "server error",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.ERROR,
			Error:   err,
		})
		os.Exit(1)
	case <-sigChan:
		logger.Log(logging.Logger{
			Context: ctx,
			Code:    200,
			Message: "signal received",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.INFO,
			Error:   nil,
		})

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log(logging.Logger{
				Context: ctx,
				Code:    500,
				Message: "shutdown error",
				From:    "main",
				Layer:   logging.LoggerLayers.SERVER,
				TypeLog: logging.LoggerTypes.ERROR,
				Error:   err,
			})
			os.Exit(1)
		}

		logger.Log(logging.Logger{
			Context: ctx,
			Code:    200,
			Message: "server shutdown complete",
			From:    "main",
			Layer:   logging.LoggerLayers.SERVER,
			TypeLog: logging.LoggerTypes.INFO,
		})
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

func healthHandler(_ *gorm.DB) gin.HandlerFunc {
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
// Helpers
// ─────────────────────────────────────────────────────────────────────────

func parseDuration(s string, defaultDuration time.Duration) time.Duration {
	if s == "" {
		return defaultDuration
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultDuration
	}
	return d
}
