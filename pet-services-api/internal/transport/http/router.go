package http

import (
	"github.com/gin-gonic/gin"

	domainauth "github.com/guilherme/pet-services-api/internal/domain/auth"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// NewRouter cria um router Gin com as rotas registradas.
func NewRouter(uc factory.UseCases, tokenService domainauth.TokenService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	// Auth
	authGroup := v1.Group("/auth")
	RegisterAuthRoutes(authGroup, NewAuthHandler(uc.Auth))

	// Users
	userGroup := v1.Group("/users")
	userGroup.Use(AuthMiddleware(tokenService))
	RegisterUserRoutes(userGroup, NewUserHandler(uc.User))

	return r
}
