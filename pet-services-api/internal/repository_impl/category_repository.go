package repository_impl

import (
	"context"
	"errors"
	"pet-services-api/internal/consts"
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

func (r *categoryRepository) FindByID(id string) (*entities.Category, error) {
	var model models.Category
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.CategoryNotFoundError)
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *categoryRepository) FindByName(name string) (*entities.Category, error) {
	var model models.Category
	err := r.db.First(&model, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.CategoryNotFoundError)
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *categoryRepository) ListCategoriesPaginated(ctx context.Context, name string, offset, limit int) ([]entities.Category, error) {
	var models []models.Category
	query := r.db.Model(&models)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	err := query.Offset(offset).Limit(limit).Order("name asc").Find(&models).Error
	if err != nil {
		return nil, err
	}
	categories := make([]entities.Category, 0, len(models))
	for _, m := range models {
		categories = append(categories, *m.ToEntity())
	}
	return categories, nil
}

func (r *categoryRepository) CountCategories(ctx context.Context, name string) (int, error) {
	var count int64
	query := r.db.Model(&models.Category{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
