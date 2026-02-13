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

func (r *photoRepository) CreateAndAttachToPet(petID string, photo *entities.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var model models.Photo
		model.FromEntity(photo)

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		pet := models.Pet{ID: petID}
		if err := tx.Model(&pet).Association("Photos").Append(&model); err != nil {
			return err
		}

		return nil
	})
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

func (r *photoRepository) ReplaceUserPhoto(userID string, photo *entities.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		user := models.User{ID: userID}

		var existing []models.Photo
		if err := tx.Model(&user).Association("Photos").Find(&existing); err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Photos").Clear(); err != nil {
			return err
		}
		if len(existing) > 0 {
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}

		var model models.Photo
		model.FromEntity(photo)
		if err := tx.Create(&model).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Photos").Append(&model); err != nil {
			return err
		}

		return nil
	})
}
