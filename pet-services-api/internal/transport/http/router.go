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

	domainauth "github.com/guilherme/pet-services-api/internal/domain/auth"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
	_ "github.com/guilherme/pet-services-api/internal/transport/http/docs"
)

// NewRouter cria um router Gin com as rotas registradas.
func NewRouter(uc factory.UseCases, tokenService domainauth.TokenService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	// Swagger UI (static spec placeholder)
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth
	authGroup := v1.Group("/auth")
	RegisterAuthRoutes(authGroup, NewAuthHandler(uc.Auth))

	// Users
	userGroup := v1.Group("/users")
	userGroup.Use(AuthMiddleware(tokenService))
	RegisterUserRoutes(userGroup, NewUserHandler(uc.User))

	// Providers
	providerPublic := v1.Group("/providers")
	RegisterProviderPublicRoutes(providerPublic, NewProviderHandler(uc.Provider))

	providerGroup := v1.Group("/providers")
	providerGroup.Use(AuthMiddleware(tokenService))
	RegisterProviderRoutes(providerGroup, NewProviderHandler(uc.Provider))

	// Requests
	requestGroup := v1.Group("/requests")
	requestGroup.Use(AuthMiddleware(tokenService))
	RegisterRequestRoutes(requestGroup, NewRequestHandler(uc.Request, uc.ProviderRepo))

	// Reviews
	reviewGroup := v1.Group("/reviews")
	reviewGroup.Use(AuthMiddleware(tokenService))
	RegisterReviewRoutes(reviewGroup, NewReviewHandler(uc.Review, uc.ProviderRepo))

	return r
}
