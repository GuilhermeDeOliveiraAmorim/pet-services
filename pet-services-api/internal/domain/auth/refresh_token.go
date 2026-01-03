package auth

import (
    "time"

    "github.com/google/uuid"
)

// RefreshToken representa um token de renovação persistido para revogação/rotacionamento.
type RefreshToken struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    ExpiresAt time.Time
    Revoked   bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewRefreshToken cria um novo registro de refresh token.
func NewRefreshToken(userID uuid.UUID, expiresAt time.Time) *RefreshToken {
    now := time.Now()
    return &RefreshToken{
        ID:        uuid.New(),
        UserID:    userID,
        ExpiresAt: expiresAt,
        Revoked:   false,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// Revoke marca o token como revogado.
func (rt *RefreshToken) Revoke() {
    rt.Revoked = true
    rt.UpdatedAt = time.Now()
}
