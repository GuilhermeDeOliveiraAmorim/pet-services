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

func (r *photoRepository) CreateAndAttachToService(serviceID string, photo *entities.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var model models.Photo
		model.FromEntity(photo)

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		service := models.Service{ID: serviceID}
		if err := tx.Model(&service).Association("Photos").Append(&model); err != nil {
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

func (r *photoRepository) CreateAndAttachToProvider(providerID string, photo *entities.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var model models.Photo
		model.FromEntity(photo)

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		provider := models.Provider{ID: providerID}
		if err := tx.Model(&provider).Association("Photos").Append(&model); err != nil {
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

func (r *photoRepository) DeleteServicePhoto(serviceID, photoID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		service := models.Service{ID: serviceID}
		photo := models.Photo{ID: photoID}

		if err := tx.Model(&service).Association("Photos").Delete(&photo); err != nil {
			return err
		}

		if err := tx.Delete(&photo).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *photoRepository) DeletePetPhoto(petID, photoID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		pet := models.Pet{ID: petID}
		photo := models.Photo{ID: photoID}

		if err := tx.Model(&pet).Association("Photos").Delete(&photo); err != nil {
			return err
		}

		if err := tx.Delete(&photo).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *photoRepository) CountProviderPhotos(providerID string) (int, error) {
	var count int64
	err := r.db.Table("provider_photos").
		Where("provider_id = ?", providerID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *photoRepository) CountPetPhotos(petID string) (int, error) {
	var count int64
	err := r.db.Table("pet_photos").
		Where("pet_id = ?", petID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *photoRepository) CountServicePhotos(serviceID string) (int, error) {
	var count int64
	err := r.db.Table("service_photos").
		Where("service_id = ?", serviceID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
