package auth

import (
	"time"

	"pet-services-api/internal/domain/user"

	"github.com/google/uuid"
)

// TokenPair representa access/refresh tokens e metadados.
type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
	RefreshID        uuid.UUID
}

// RefreshClaims representa os dados extraídos de um refresh token.
type RefreshClaims struct {
	TokenID   uuid.UUID
	UserID    uuid.UUID
	UserType  user.UserType
	ExpiresAt time.Time
}

// AccessClaims representa os dados extraídos de um access token.
type AccessClaims struct {
	UserID    uuid.UUID
	UserType  user.UserType
	ExpiresAt time.Time
}

// TokenService abstrai geração e validação de tokens (JWT, etc.).
type TokenService interface {
	GenerateTokens(userID uuid.UUID, userType user.UserType) (TokenPair, error)
	ParseRefreshToken(token string) (RefreshClaims, error)
	ParseAccessToken(token string) (AccessClaims, error)
}
