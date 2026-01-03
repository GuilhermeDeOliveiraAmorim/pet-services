package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken armazena tokens de renovação.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_refresh_tokens_user"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela no banco.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// EmailVerificationToken guarda tokens de verificação de email.
type EmailVerificationToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_email_tokens_user"`
	Token     string    `gorm:"size:255;not null;uniqueIndex:idx_email_tokens_token"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName define o nome da tabela no banco.
func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// PasswordResetToken guarda tokens de reset de senha.
type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_password_reset_tokens_user"`
	Token     string    `gorm:"size:255;not null;uniqueIndex:idx_password_reset_tokens_token"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName define o nome da tabela no banco.
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
