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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c)
		if tokenString == "" {
			return
		}

		token, claims := parseToken(c, tokenString)
		if token == nil || claims == nil {
			return
		}

		userType, okType := claims["user_type"].(string)
		userID, okID := claims["id"].(string)
		if !okType || !okID {
			abortUnauthorized(c, "Token inválido", "Claims incompletas")
			return
		}

		c.Set("user_id", userID)
		c.Set("user_type", userType)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortUnauthorized(c, "Autorização não fornecida", "Cabeçalho de autorização está vazio")
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		abortUnauthorized(c, "Bearer token inválido", "Cabeçalho de autorização inválido")
		return ""
	}

	return parts[1]
}

func parseToken(c *gin.Context, tokenString string) (*jwt.Token, jwt.MapClaims) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			abortUnauthorized(c, "Token inválido", "Método de assinatura inesperado")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.SECRETS_VAR.JWT_SECRET), nil
	})

	if err != nil || token == nil || !token.Valid {
		abortUnauthorized(c, "Token inválido", "Erro ao analisar o token")
		return nil, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		abortUnauthorized(c, "Token inválido", "Não foi possível ler as claims")
		return nil, nil
	}

	return token, claims
}

func abortUnauthorized(c *gin.Context, title, detail string) {
	problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
		Title:  title,
		Detail: detail,
	})

	logging.NewLogger(logging.Logger{
		Context:  c.Request.Context(),
		Code:     exceptions.RFC401_CODE,
		Message:  title,
		From:     "AuthMiddleware",
		Layer:    logging.LoggerLayers.MIDDLEWARES,
		TypeLog:  logging.LoggerTypes.ERROR,
		Error:    errors.New(detail),
		Problems: []exceptions.ProblemDetails{problem},
	})

	c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
}

func OwnerOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "owner" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo owner",
			})
			logging.NewLogger(logging.Logger{
				Context:  c.Request.Context(),
				Code:     exceptions.RFC401_CODE,
				Message:  problem.Title,
				From:     "AuthMiddleware",
				Layer:    logging.LoggerLayers.MIDDLEWARES,
				TypeLog:  logging.LoggerTypes.ERROR,
				Error:    errors.New(problem.Detail),
				Problems: []exceptions.ProblemDetails{problem},
			})
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func ProviderOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "provider" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo provider",
			})
			logging.NewLogger(logging.Logger{
				Context:  c.Request.Context(),
				Code:     exceptions.RFC401_CODE,
				Message:  problem.Title,
				From:     "AuthMiddleware",
				Layer:    logging.LoggerLayers.MIDDLEWARES,
				TypeLog:  logging.LoggerTypes.ERROR,
				Error:    errors.New(problem.Detail),
				Problems: []exceptions.ProblemDetails{problem},
			})
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para administradores",
			})
			logging.NewLogger(logging.Logger{
				Context:  c.Request.Context(),
				Code:     exceptions.RFC401_CODE,
				Message:  problem.Title,
				From:     "AuthMiddleware",
				Layer:    logging.LoggerLayers.MIDDLEWARES,
				TypeLog:  logging.LoggerTypes.ERROR,
				Error:    errors.New(problem.Detail),
				Problems: []exceptions.ProblemDetails{problem},
			})
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}
