package repository_impl

import (
	"context"
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

func (r *tagRepository) ListTagsPaginated(ctx context.Context, name string, offset, limit int) ([]entities.Tag, error) {
	var tagModels []models.Tag
	query := r.db.Model(&models.Tag{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	err := query.Offset(offset).Limit(limit).Order("name asc").Find(&tagModels).Error
	if err != nil {
		return nil, err
	}
	tags := make([]entities.Tag, 0, len(tagModels))
	for _, m := range tagModels {
		tags = append(tags, *m.ToEntity())
	}
	return tags, nil
}

func (r *tagRepository) CountTags(ctx context.Context, name string) (int, error) {
	var count int64
	query := r.db.Model(&models.Tag{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
