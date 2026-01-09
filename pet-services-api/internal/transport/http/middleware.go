package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

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

		// Injeta no contexto Gin para uso dos handlers
		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxUserTypeKey, claims.UserType)

		c.Next()
	}
}
