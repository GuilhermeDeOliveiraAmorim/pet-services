package auth

import (
    "time"

    "github.com/google/uuid"
    "github.com/guilherme/pet-services-api/internal/domain/user"
)

// TokenPair representa access/refresh tokens e metadados.
type TokenPair struct {
    AccessToken     string
    AccessExpiresAt time.Time
    RefreshToken    string
    RefreshExpiresAt time.Time
    RefreshID       uuid.UUID
}

// RefreshClaims representa os dados extraídos de um refresh token.
type RefreshClaims struct {
    TokenID   uuid.UUID
    UserID    uuid.UUID
    UserType  user.UserType
    ExpiresAt time.Time
}

// TokenService abstrai geração e validação de tokens (JWT, etc.).
type TokenService interface {
    GenerateTokens(userID uuid.UUID, userType user.UserType) (TokenPair, error)
    ParseRefreshToken(token string) (RefreshClaims, error)
}
