package repository_impl

import (
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) entities.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *entities.Service) error {
	var model models.Service
	model.FromEntity(service)
	return r.db.Create(&model).Error
}

func (r *serviceRepository) Update(service *entities.Service) error {
	var model models.Service
	model.FromEntity(service)
	return r.db.Save(&model).Error
}

func (r *serviceRepository) FindByID(id string) (*entities.Service, error) {
	var model models.Service
	err := r.db.
		Preload("Photos").
		Preload("Categories").
		Preload("Tags").
		First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.ServiceNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *serviceRepository) HasTag(serviceID, tagID string) (bool, error) {
	var count int64
	err := r.db.Table("service_tags").
		Where("service_id = ? AND tag_id = ?", serviceID, tagID).
		Count(&count).Error
	return count > 0, err
}

func (r *serviceRepository) AddTag(serviceID, tagID string) error {
	service := models.Service{ID: serviceID}
	tag := models.Tag{ID: tagID}
	return r.db.Model(&service).Association("Tags").Append(&tag)
}

func (r *serviceRepository) HasCategory(serviceID, categoryID string) (bool, error) {
	var count int64
	err := r.db.Table("service_categories").
		Where("service_id = ? AND category_id = ?", serviceID, categoryID).
		Count(&count).Error
	return count > 0, err
}

func (r *serviceRepository) AddCategory(serviceID, categoryID string) error {
	service := models.Service{ID: serviceID}
	category := models.Category{ID: categoryID}
	return r.db.Model(&service).Association("Categories").Append(&category)
}

func (r *serviceRepository) RemoveCategory(serviceID, categoryID string) error {
	service := models.Service{ID: serviceID}
	category := models.Category{ID: categoryID}
	return r.db.Model(&service).Association("Categories").Delete(&category)
}

func (r *serviceRepository) List(providerID, categoryID, tagID string, priceMin, priceMax float64, page, pageSize int) ([]*entities.Service, int64, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).
		Preload("Provider").
		Preload("Photos").
		Preload("Categories").
		Preload("Tags").
		Where("active = ?", true)

	if providerID != "" {
		query = query.Where("provider_id = ?", providerID)
	}

	if categoryID != "" {
		query = query.Joins("INNER JOIN service_categories ON service_categories.service_id = services.id").
			Where("service_categories.category_id = ?", categoryID)
	}

	if tagID != "" {
		query = query.Joins("INNER JOIN service_tags ON service_tags.service_id = services.id").
			Where("service_tags.tag_id = ?", tagID)
	}

	if priceMin > 0 {
		query = query.Where("price_minimum >= ?", priceMin)
	}

	if priceMax > 0 {
		query = query.Where("price_maximum <= ?", priceMax)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entities.Service, len(services))
	for i, svc := range services {
		entities[i] = svc.ToEntity()
	}

	return entities, total, nil
}

func (r *serviceRepository) Search(query, categoryID, tagID string, latitude, longitude, radiusKm, priceMin, priceMax float64, page, pageSize int) ([]*entities.Service, int64, error) {
	var services []models.Service
	var total int64

	baseQuery := r.db.Model(&models.Service{}).
		Preload("Provider").
		Preload("Photos").
		Preload("Categories").
		Preload("Tags").
		Where("services.active = ?", true)

	// Busca textual por nome ou descrição
	if query != "" {
		searchPattern := "%" + query + "%"
		baseQuery = baseQuery.Where("services.name ILIKE ? OR services.description ILIKE ?", searchPattern, searchPattern)
	}

	// Filtro por categoria
	if categoryID != "" {
		baseQuery = baseQuery.Joins("INNER JOIN service_categories ON service_categories.service_id = services.id").
			Where("service_categories.category_id = ?", categoryID)
	}

	// Filtro por tag
	if tagID != "" {
		baseQuery = baseQuery.Joins("INNER JOIN service_tags ON service_tags.service_id = services.id").
			Where("service_tags.tag_id = ?", tagID)
	}

	// Busca geoespacial por raio
	if latitude != 0 && longitude != 0 && radiusKm > 0 {
		// Fórmula Haversine para calcular distância em km
		// earthRadiusKm = 6371
		baseQuery = baseQuery.Joins("INNER JOIN providers ON providers.id = services.provider_id").
			Where("providers.active = ?", true).
			Where("(6371 * acos(cos(radians(?)) * cos(radians(providers.latitude)) * cos(radians(providers.longitude) - radians(?)) + sin(radians(?)) * sin(radians(providers.latitude)))) <= ?",
				latitude, longitude, latitude, radiusKm)
	} else {
		// Apenas garantir que o provider está ativo
		baseQuery = baseQuery.Joins("INNER JOIN providers ON providers.id = services.provider_id").
			Where("providers.active = ?", true)
	}

	// Filtro por faixa de preço
	if priceMin > 0 {
		baseQuery = baseQuery.Where("services.price_minimum >= ?", priceMin)
	}

	if priceMax > 0 {
		baseQuery = baseQuery.Where("services.price_maximum <= ?", priceMax)
	}

	// Contar total de resultados
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginação e ordenação
	offset := (page - 1) * pageSize
	if err := baseQuery.Order("services.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entities.Service, len(services))
	for i, svc := range services {
		entities[i] = svc.ToEntity()
	}

	return entities, total, nil
}
