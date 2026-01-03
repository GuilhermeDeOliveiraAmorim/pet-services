package auth

import (
    "context"

    "github.com/google/uuid"
)

// RefreshTokenRepository define persistência e revogação de refresh tokens.
type RefreshTokenRepository interface {
    Create(ctx context.Context, token *RefreshToken) error
    FindByID(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
    Revoke(ctx context.Context, id uuid.UUID) error
    RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
