package middlewares

import (
	"errors"
	"net/http"
	"pet-services-api/internal/config"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c, logger)
		if tokenString == "" {
			return
		}

		token, claims := parseToken(c, tokenString, logger)
		if token == nil || claims == nil {
			return
		}

		userType, okType := claims["user_type"].(string)
		userID, okID := claims["id"].(string)
		if !okType || !okID {
			abortUnauthorized(c, "Token inválido", "Claims incompletas", logger)
			return
		}

		c.Set("user_id", userID)
		c.Set("user_type", userType)
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

func parseToken(c *gin.Context, tokenString string, logger logging.LoggerInterface) (*jwt.Token, jwt.MapClaims) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			abortUnauthorized(c, "Token inválido", "Método de assinatura inesperado", logger)
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.SECRETS_VAR.JWT_SECRET), nil
	})

	if err != nil || token == nil || !token.Valid {
		abortUnauthorized(c, "Token inválido", "Erro ao analisar o token", logger)
		return nil, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		abortUnauthorized(c, "Token inválido", "Não foi possível ler as claims", logger)
		return nil, nil
	}

	return token, claims
}

func abortUnauthorized(c *gin.Context, title, detail string, logger logging.LoggerInterface) {
	problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
		Title:  title,
		Detail: detail,
	})

	logger.LogError(c.Request.Context(), "AuthMiddleware", title, errors.New(detail))

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
