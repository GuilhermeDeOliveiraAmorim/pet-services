package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) entities.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entities.Category) error {
	var model models.Category
	model.FromEntity(category)
	return r.db.Create(&model).Error
}

func (r *categoryRepository) FindByName(name string) (*entities.Category, error) {
	var model models.Category
	err := r.db.First(&model, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}
