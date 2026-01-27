package repository_impl

import (
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entities.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
	var model models.User
	model.FromEntity(user)
	return r.db.Create(&model).Error
}

func (r *userRepository) FindByID(id string) (*entities.User, error) {
	var model models.User
	err := r.db.Preload("Photos").Preload("Pets").First(&model, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("usuario não encontrado")
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var model models.User
	err := r.db.Preload("Photos").Preload("Pets").First(&model, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("usuario não encontrado")
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByPhone(countryCode, areaCode, number string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("country_code = ? AND area_code = ? AND phone_number = ?", countryCode, areaCode, number).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) UpdateEmailVerified(userID string, verified bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("email_verified", verified).Error
}

func (r *userRepository) Update(user *entities.User) error {
	var model models.User
	model.FromEntity(user)
	return r.db.Save(&model).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) List(page, limit int) ([]*entities.User, int64, error) {
	var modelsUsers []models.User
	var total int64
	offset := (page - 1) * limit
	db := r.db.Model(&models.User{})
	db.Count(&total)
	err := db.Offset(offset).Limit(limit).Preload("Photos").Preload("Pets").Find(&modelsUsers).Error
	if err != nil {
		return nil, 0, err
	}
	var users []*entities.User
	for _, m := range modelsUsers {
		users = append(users, m.ToEntity())
	}
	return users, total, nil
}
