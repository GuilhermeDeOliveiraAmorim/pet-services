// Package http expõe rotas HTTP e middlewares.
//
// @title Pet Services API
// @version 1.0.0
// @description API para autenticação, usuários, prestadores, solicitações e avaliações.
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	domainauth "pet-services-api/internal/domain/auth"
	"pet-services-api/internal/infra/factory"
	_ "pet-services-api/internal/transport/http/docs"
)

// NewRouter cria um router Gin com as rotas registradas.
func NewRouter(uc factory.UseCases, tokenService domainauth.TokenService) *gin.Engine {
	r := gin.Default()

	// Global rate limiting middleware
	globalRateLimiter := NewRateLimiter()
	r.Use(globalRateLimiter.RateLimitMiddleware())

	v1 := r.Group("/api/v1")

	// Swagger UI com suporte a Bearer Token
	swaggerConfig := &ginSwagger.Config{
		URL:                  "/api/v1/swagger/doc.json", // Caminho do spec
		DeepLinking:          true,
		PersistAuthorization: true, // Mantém o Bearer Token
	}

	v1.GET("/swagger/*any", ginSwagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))

	errorService := NewErrorService()

	// Auth (stricter rate limiting for security)
	authGroup := v1.Group("/auth")
	authGroup.Use(CriticalEndpointRateLimiter())
	RegisterAuthRoutes(authGroup, NewAuthHandler(uc.Auth, errorService))

	// Users
	userGroup := v1.Group("/users")
	userGroup.Use(AuthMiddleware(tokenService))
	RegisterUserRoutes(userGroup, NewUserHandler(uc.User, errorService))

	// Providers
	v1.GET("/providers/search/location", NewProviderHandler(uc.Provider, errorService).ListByLocation)

	// Grupo autenticado: todas rotas protegidas
	providerGroup := v1.Group("/providers")
	providerGroup.Use(AuthMiddleware(tokenService))
	RegisterProviderRoutes(providerGroup, NewProviderHandler(uc.Provider, errorService))

	// Requests
	requestGroup := v1.Group("/requests")
	requestGroup.Use(AuthMiddleware(tokenService))
	RegisterRequestRoutes(requestGroup, NewRequestHandler(uc.Request, uc.ProviderRepo, errorService))

	// Reviews
	reviewGroup := v1.Group("/reviews")
	reviewGroup.Use(AuthMiddleware(tokenService))
	RegisterReviewRoutes(reviewGroup, NewReviewHandler(uc.Review, uc.ProviderRepo, errorService))

	return r
}
