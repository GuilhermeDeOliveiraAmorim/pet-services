package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken armazena tokens de renovação.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_refresh_tokens_user;index:idx_refresh_tokens_user_revoked,composite:idx_refresh_tokens_user_revoked;foreignKey:UserID;references:ID"`
	ExpiresAt time.Time `gorm:"not null;index:idx_refresh_tokens_expires"`
	Revoked   bool      `gorm:"default:false;index:idx_refresh_tokens_user_revoked,composite:idx_refresh_tokens_user_revoked"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName define o nome da tabela no banco.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// EmailVerificationToken guarda tokens de verificação de email.
type EmailVerificationToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_email_tokens_user;index:idx_email_tokens_user_used,composite:idx_email_tokens_user_used;foreignKey:UserID;references:ID"`
	Token     string    `gorm:"size:255;not null;uniqueIndex:idx_email_tokens_token"`
	ExpiresAt time.Time `gorm:"not null;index:idx_email_tokens_expires"`
	Used      bool      `gorm:"default:false;index:idx_email_tokens_user_used,composite:idx_email_tokens_user_used"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName define o nome da tabela no banco.
func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// PasswordResetToken guarda tokens de reset de senha.
type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_password_reset_tokens_user;index:idx_password_reset_tokens_user_used,composite:idx_password_reset_tokens_user_used;foreignKey:UserID;references:ID"`
	Token     string    `gorm:"size:255;not null;uniqueIndex:idx_password_reset_tokens_token"`
	ExpiresAt time.Time `gorm:"not null;index:idx_password_reset_tokens_expires"`
	Used      bool      `gorm:"default:false;index:idx_password_reset_tokens_user_used,composite:idx_password_reset_tokens_user_used"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName define o nome da tabela no banco.
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
