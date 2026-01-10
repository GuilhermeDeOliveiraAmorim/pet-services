package http

import (
	"errors"
	"log/slog"
	"net/http"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainauth "pet-services-api/internal/domain/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Extrai o token Bearer do header Authorization, abortando com 401 se inválido.
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	slog.Info("TOKEN", slog.String("auth_header", authHeader))

	if authHeader == "" {
		c.AbortWithStatusJSON(401, errorResponse("missing_authorization", "Cabeçalho de autorização está vazio"))
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(401, errorResponse("invalid_authorization", "Bearer token inválido"))
		return ""
	}

	slog.Info("TOKEN", slog.String("part1", parts[0]), slog.String("part2", parts[1]))

	token := strings.TrimSpace(parts[1])
	if token == "" {
		c.AbortWithStatusJSON(401, errorResponse("invalid_authorization", "Token está vazio"))
		return ""
	}

	slog.Info("TOKEN", slog.String("token", token))

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

// parseToken agora recebe o segredo JWT como argumento.
func parseToken(c *gin.Context, tokenString string, jwtSecret string) (*jwt.Token, jwt.MapClaims) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			abortUnauthorized(c, "Token inválido", "Método de assinatura inesperado")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
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

	// Loga o erro de autenticação usando slog padrão
	slog.Error("AuthMiddleware unauthorized",
		slog.String("title", title),
		slog.String("detail", detail),
		slog.Int("code", exceptions.RFC401_CODE),
		slog.String("layer", logging.LoggerLayers.MIDDLEWARES),
	)

	c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
}
