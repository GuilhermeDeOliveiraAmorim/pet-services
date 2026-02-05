package middlewares

import (
	"errors"
	"net/http"
	"pet-services-api/internal/auth"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c, logger)
		if tokenString == "" {
			return
		}

		jwtService, err := auth.NewJWTServiceFromEnv()
		if err != nil {
			abortUnauthorized(c, "Configuração JWT inválida", "Erro ao carregar configuração de autenticação", logger)
			return
		}

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			abortUnauthorized(c, "Token inválido", err.Error(), logger)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context, logger logging.LoggerInterface) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortUnauthorized(c, "Autorização não fornecida", "Cabeçalho de autorização está vazio", logger)
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		abortUnauthorized(c, "Bearer token inválido", "Cabeçalho de autorização inválido", logger)
		return ""
	}

	return parts[1]
}

func abortUnauthorized(c *gin.Context, title, detail string, logger logging.LoggerInterface) {
	problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
		Title:  title,
		Detail: detail,
	})

	logger.LogError(c.Request.Context(), "AuthMiddleware", title, nil)

	c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
}

func OwnerOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "owner" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo owner",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func ProviderOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "provider" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo provider",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func AdminOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para administradores",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}
