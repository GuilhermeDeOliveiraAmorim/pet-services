package routes

import (
	"context"
	"pet-services-api/internal/config"
	"pet-services-api/internal/database"
	"pet-services-api/internal/handlers"
	"pet-services-api/internal/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(storageInput database.StorageInput, ctx context.Context) *gin.Engine {
	handlerFactory := handlers.NewHandlerFactory(storageInput)
	middlewareFactory := middlewares.NewMiddlewareFactory()

	r := gin.Default()

	devURL, prodURL := config.GetFrontendURLs()

	allowOrigins := []string{}
	if devURL != "" {
		allowOrigins = append(allowOrigins, devURL)
	}
	if prodURL != "" {
		allowOrigins = append(allowOrigins, prodURL)
	}
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewareFactory.DefaultRateLimitMiddleware())

	public := r.Group("/")
	{
		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Health check endpoints
		// public.GET("/health", handlerFactory.HealthHandler.Check)
		// public.GET("/health/ready", handlerFactory.HealthHandler.Ready)
		// public.GET("/health/live", handlerFactory.HealthHandler.Live)

		public.POST("/users/register", handlerFactory.UserHandler.RegisterUser)
		public.POST("/users/check-email", handlerFactory.UserHandler.CheckEmailExists)
		public.POST("/users/check-phone", handlerFactory.UserHandler.CheckPhoneExists)

		public.POST("/auth/login", handlerFactory.TokenHandler.LoginUser)
		public.POST("/auth/request-password-reset", handlerFactory.TokenHandler.RequestPasswordReset)
		public.POST("/auth/reset-password", handlerFactory.TokenHandler.ResetPassword)
		public.POST("/auth/resend-verification-email", handlerFactory.TokenHandler.ResendVerificationEmail)
		public.POST("/auth/verify-email", handlerFactory.TokenHandler.VerifyEmail)
		public.POST("/auth/logout", handlerFactory.TokenHandler.Logout)
	}

	authorized := r.Group("/")
	authorized.Use(middlewareFactory.AuthMiddleware())
	{
		authorized.GET("/users/profile", handlerFactory.UserHandler.GetProfile)
		authorized.GET("/users/:user_id", handlerFactory.UserHandler.GetUserByID)
		authorized.GET("/users", handlerFactory.UserHandler.ListUsers)
		authorized.PUT("/users", handlerFactory.UserHandler.UpdateUser)
		authorized.DELETE("/users", handlerFactory.UserHandler.DeleteUser)
		authorized.POST("/users/reactivate", handlerFactory.UserHandler.ReactivateUser)
		authorized.POST("/users/deactivate", handlerFactory.UserHandler.DeactivateUser)
		authorized.POST("/users/change-password", handlerFactory.UserHandler.ChangePassword)
		authorized.POST("/users/update-email-verified", handlerFactory.UserHandler.UpdateEmailVerified)
	}

	return r
}
