package user

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken representa um token de reset de senha.
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// NewPasswordResetToken cria um novo token de reset com validade.
func NewPasswordResetToken(userID uuid.UUID, token string, expiresAt time.Time) *PasswordResetToken {
	return &PasswordResetToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}
}

// IsValid verifica se o token ainda é válido.
func (t *PasswordResetToken) IsValid() bool {
	return !t.Used && time.Now().Before(t.ExpiresAt)
}

// MarkAsUsed marca o token como já utilizado.
func (t *PasswordResetToken) MarkAsUsed() {
	t.Used = true
}
