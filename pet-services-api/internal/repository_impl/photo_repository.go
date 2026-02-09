package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) entities.PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) CreateAndAttachToUser(userID string, photo *entities.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var model models.Photo
		model.FromEntity(photo)

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		user := models.User{ID: userID}
		if err := tx.Model(&user).Association("Photos").Append(&model); err != nil {
			return err
		}

		return nil
	})
}
