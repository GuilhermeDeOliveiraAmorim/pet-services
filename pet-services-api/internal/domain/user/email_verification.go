package user

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerificationToken representa um token de verificação de email.
type EmailVerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// NewEmailVerificationToken cria um novo token de verificação com validade.
func NewEmailVerificationToken(userID uuid.UUID, token string, expiresAt time.Time) *EmailVerificationToken {
	return &EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}
}

// IsValid verifica se o token ainda é válido.
func (t *EmailVerificationToken) IsValid() bool {
	return !t.Used && time.Now().Before(t.ExpiresAt)
}

// MarkAsUsed marca o token como já utilizado.
func (t *EmailVerificationToken) MarkAsUsed() {
	t.Used = true
}
