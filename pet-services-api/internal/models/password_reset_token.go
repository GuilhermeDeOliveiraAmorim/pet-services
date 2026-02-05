package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID    string         `gorm:"type:char(26);not null;index:idx_user_active,priority:1" json:"user_id"`
	Token     string         `gorm:"type:varchar(500);uniqueIndex;not null" json:"token"`
	TokenType string         `gorm:"type:varchar(20);not null;default:'reset';index" json:"token_type"`
	ExpiresAt time.Time      `gorm:"not null;index:idx_expires" json:"expires_at"`
	RevokedAt *time.Time     `json:"revoked_at"`
	UserAgent string         `gorm:"type:varchar(255)" json:"user_agent"`
	IpAddress string         `gorm:"type:varchar(45)" json:"ip_address"`
	Active    bool           `gorm:"default:true;index:idx_user_active,priority:2" json:"active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

func (prt *PasswordResetToken) ToEntity() *entities.PasswordResetToken {
	return &entities.PasswordResetToken{
		Token:     prt.Token,
		UserID:    prt.UserID,
		ExpiresAt: prt.ExpiresAt,
		UserAgent: prt.UserAgent,
		IP:        prt.IpAddress,
		RevokedAt: prt.RevokedAt,
	}
}

func (prt *PasswordResetToken) FromEntity(entity *entities.PasswordResetToken) {
	prt.ID = entity.Token
	prt.UserID = entity.UserID
	prt.Token = entity.Token
	prt.ExpiresAt = entity.ExpiresAt
	prt.UserAgent = entity.UserAgent
	prt.IpAddress = entity.IP
	prt.RevokedAt = entity.RevokedAt
	prt.TokenType = "reset"
	prt.Active = entity.RevokedAt == nil
}
