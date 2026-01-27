package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Services []Service `gorm:"many2many:service_categories" json:"services"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) ToEntity() *entities.Category {
	return &entities.Category{
		Base: entities.Base{
			ID:            c.ID,
			Active:        c.Active,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
			DeactivatedAt: c.DeactivatedAt,
		},
		Name: c.Name,
	}
}

func (c *Category) FromEntity(entity *entities.Category) {
	c.ID = entity.ID
	c.Name = entity.Name
	c.Active = entity.Active
	c.CreatedAt = entity.CreatedAt
	c.UpdatedAt = entity.UpdatedAt
	c.DeactivatedAt = entity.DeactivatedAt
}
