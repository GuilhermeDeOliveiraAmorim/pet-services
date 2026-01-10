package http

import (
	"log/slog"
	"net/http"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainauth "pet-services-api/internal/domain/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Extrai o token Bearer do header Authorization, abortando com 401 se inválido.
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(401, errorResponse("missing_authorization", "Cabeçalho de autorização está vazio"))
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(401, errorResponse("invalid_authorization", "Bearer token inválido"))
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		c.AbortWithStatusJSON(401, errorResponse("invalid_authorization", "Token está vazio"))
		return ""
	}

	return token
}

const (
	ctxUserIDKey   = "user_id"
	ctxUserTypeKey = "user_type"
)

// AuthMiddleware extrai dados do JWT (usa refresh token parsing disponível) e injeta user_id/user_type no contexto.
// Espera header Authorization: Bearer <token>.
func AuthMiddleware(tokenService domainauth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			return
		}

		claims, err := tokenService.ParseAccessToken(token)
		if err != nil {
			abortUnauthorized(c, "Token inválido", err.Error())
			return
		}

		userType := claims.UserType
		userID := claims.UserID
		if userID == uuid.Nil {
			abortUnauthorized(c, "Token inválido", "user_id ausente nas claims")
			return
		}
		c.Set(ctxUserIDKey, userID.String())
		c.Set(ctxUserTypeKey, string(userType))

		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, title, detail string) {
	problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
		Title:  title,
		Detail: detail,
	})

	// Loga o erro de autenticação usando slog padrão
	slog.Error("AuthMiddleware unauthorized",
		slog.String("title", title),
		slog.String("detail", detail),
		slog.Int("code", exceptions.RFC401_CODE),
		slog.String("layer", logging.LoggerLayers.MIDDLEWARES),
	)

	c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
}
