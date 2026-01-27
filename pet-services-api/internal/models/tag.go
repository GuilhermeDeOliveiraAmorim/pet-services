package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(30);uniqueIndex;not null" json:"name"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Services []Service `gorm:"many2many:service_tags" json:"services"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) ToEntity() *entities.Tag {
	return &entities.Tag{
		Base: entities.Base{
			ID:            t.ID,
			Active:        t.Active,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
			DeactivatedAt: t.DeactivatedAt,
		},
		Name: t.Name,
	}
}

func (t *Tag) FromEntity(entity *entities.Tag) {
	t.ID = entity.ID
	t.Name = entity.Name
	t.Active = entity.Active
	t.CreatedAt = entity.CreatedAt
	t.UpdatedAt = entity.UpdatedAt
	t.DeactivatedAt = entity.DeactivatedAt
}
