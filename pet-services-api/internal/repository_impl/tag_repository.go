package repository_impl

import (
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) entities.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *entities.Tag) error {
	var model models.Tag
	model.FromEntity(tag)
	return r.db.Create(&model).Error
}

func (r *tagRepository) FindByID(id string) (*entities.Tag, error) {
	var model models.Tag
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.TagNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *tagRepository) FindByName(name string) (*entities.Tag, error) {
	var model models.Tag
	err := r.db.First(&model, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.TagNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}