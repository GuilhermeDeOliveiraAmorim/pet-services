package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	domainauth "pet-services-api/internal/domain/auth"
)

const (
	ctxUserIDKey   = "user_id"
	ctxUserTypeKey = "user_type"
)

// AuthMiddleware extrai dados do JWT (usa refresh token parsing disponível) e injeta user_id/user_type no contexto.
// Espera header Authorization: Bearer <token>.
func AuthMiddleware(tokenService domainauth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("missing_authorization", "authorization header is required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid_authorization", "expected Bearer token"))
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid_authorization", "token is empty"))
			return
		}

		claims, err := tokenService.ParseRefreshToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid_token", err.Error()))
			return
		}

		

		// Garante que claims.UserID é uuid.UUID
		var userID uuid.UUID
		switch v := any(claims.UserID).(type) {
		case uuid.UUID:
			userID = v
		case string:
			parsed, err := uuid.Parse(v)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid_token", "user_id inválido no token"))
				return
			}
			userID = parsed
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid_token", "user_id ausente ou inválido no token"))
			return
		}

		c.Set(ctxUserIDKey, userID)
		c.Set(ctxUserTypeKey, claims.UserType)

		c.Next()
	}
}
