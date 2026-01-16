package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index:idx_user_active,priority:1" json:"user_id"`
	Token         string         `gorm:"type:varchar(500);uniqueIndex;not null" json:"token"`
	ExpiresAt     time.Time      `gorm:"not null;index:idx_expires" json:"expires_at"`
	RevokedAt     *time.Time     `json:"revoked_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	UserAgent     string         `gorm:"type:varchar(255)" json:"user_agent"`
	IpAddress     string         `gorm:"type:varchar(45)" json:"ip_address"`
	Active        bool           `gorm:"default:true;index:idx_user_active,priority:2" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt *RefreshToken) ToEntity() *entities.RefreshToken {
	return &entities.RefreshToken{
		Base: entities.Base{
			ID:            rt.ID,
			Active:        rt.Active,
			CreatedAt:     rt.CreatedAt,
			UpdatedAt:     rt.UpdatedAt,
			DeactivatedAt: rt.DeactivatedAt,
		},
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: rt.RevokedAt,
		UserAgent: rt.UserAgent,
		IpAddress: rt.IpAddress,
	}
}

func (rt *RefreshToken) FromEntity(entity *entities.RefreshToken) {
	rt.ID = entity.ID
	rt.Active = entity.Active
	rt.CreatedAt = entity.CreatedAt
	rt.UpdatedAt = entity.UpdatedAt
	rt.DeactivatedAt = entity.DeactivatedAt
	rt.UserID = entity.UserID
	rt.Token = entity.Token
	rt.ExpiresAt = entity.ExpiresAt
	rt.RevokedAt = entity.RevokedAt
	rt.UserAgent = entity.UserAgent
	rt.IpAddress = entity.IpAddress
}
