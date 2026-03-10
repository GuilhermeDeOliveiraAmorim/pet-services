package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type Service struct {
	Base
	ProviderID   string     `json:"provider_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        float64    `json:"price"`
	PriceMinimum float64    `json:"price_minimum"`
	PriceMaximum float64    `json:"price_maximum"`
	Duration     int        `json:"duration"`
	Photos       []Photo    `json:"photos"`
	Categories   []Category `json:"categories"`
	Tags         []Tag      `json:"tags"`
}

type ServiceRepository interface {
	Create(service *Service) error
	FindByID(id string) (*Service, error)
	Update(service *Service) error
	List(providerID, categoryID, tagID string, priceMin, priceMax float64, page, pageSize int) ([]*Service, int64, error)
	Search(query, categoryID, tagID string, latitude, longitude, radiusKm, priceMin, priceMax float64, page, pageSize int) ([]*Service, int64, error)
	HasTag(serviceID, tagID string) (bool, error)
	AddTag(serviceID, tagID string) error
	HasCategory(serviceID, categoryID string) (bool, error)
	AddCategory(serviceID, categoryID string) error
	RemoveCategory(serviceID, categoryID string) error
}

func NewService(providerID string, name string, description string, price float64, priceMinimum float64, priceMaximum float64, duration int) (*Service, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if providerID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do provedor ausente",
			Detail: "O ID do provedor é obrigatório",
		}))
	}

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome do serviço ausente",
			Detail: "O nome do serviço é obrigatório",
		}))
	} else if len(name) > 100 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome do serviço deve ter no máximo 100 caracteres",
		}))
	}

	if description == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição do serviço ausente",
			Detail: "A descrição do serviço é obrigatória",
		}))
	} else if len(description) > 1000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição muito longa",
			Detail: "A descrição deve ter no máximo 1000 caracteres",
		}))
	}

	if price < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Preço do serviço inválido",
			Detail: "O preço do serviço não pode ser negativo",
		}))
	}

	if priceMinimum < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Preço mínimo do serviço inválido",
			Detail: "O preço mínimo do serviço não pode ser negativo",
		}))
	}

	if priceMaximum < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Preço máximo do serviço inválido",
			Detail: "O preço máximo do serviço não pode ser negativo",
		}))
	}

	if price > 0 && (priceMinimum > 0 || priceMaximum > 0) {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Conflito de preços",
			Detail: "Use 'price' para preço fixo OU 'price_minimum' e 'price_maximum' para faixa de preço, não ambos",
		}))
	}

	if priceMinimum > 0 && priceMaximum > 0 && priceMinimum > priceMaximum {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Faixa de preço inválida",
			Detail: "O preço mínimo não pode ser maior que o preço máximo",
		}))
	}

	if price == 0 && priceMinimum == 0 && priceMaximum == 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Preço ausente",
			Detail: "Defina um preço fixo ou uma faixa de preço para o serviço",
		}))
	}

	if duration < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Duração do serviço inválida",
			Detail: "A duração do serviço não pode ser negativa",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Service{
		Base:         *NewBase(),
		ProviderID:   providerID,
		Name:         name,
		Description:  description,
		Price:        price,
		PriceMinimum: priceMinimum,
		PriceMaximum: priceMaximum,
		Duration:     duration,
	}, nil
}

func (s *Service) AddPhoto(photo Photo) {
	timeNow := time.Now()
	s.Photos = append(s.Photos, photo)
	s.UpdatedAt = &timeNow
}

func (s *Service) AddCategory(category Category) {
	timeNow := time.Now()
	s.Categories = append(s.Categories, category)
	s.UpdatedAt = &timeNow
}

func (s *Service) AddTag(tag Tag) {
	timeNow := time.Now()
	s.Tags = append(s.Tags, tag)
	s.UpdatedAt = &timeNow
}

func (s *Service) RemovePhoto(photo Photo) {
	timeNow := time.Now()
	var updatedPhotos []Photo
	for _, pht := range s.Photos {
		if pht.ID != photo.ID {
			updatedPhotos = append(updatedPhotos, pht)
		}
	}
	s.Photos = updatedPhotos
	s.UpdatedAt = &timeNow
}

func (s *Service) RemoveCategory(categoryID string) {
	for i, category := range s.Categories {
		if category.ID == categoryID {
			timeNow := time.Now()
			s.Categories = append(s.Categories[:i], s.Categories[i+1:]...)
			s.UpdatedAt = &timeNow
			return
		}
	}
}

func (s *Service) RemoveTag(tagID string) {
	for i, tag := range s.Tags {
		if tag.ID == tagID {
			timeNow := time.Now()
			s.Tags = append(s.Tags[:i], s.Tags[i+1:]...)
			s.UpdatedAt = &timeNow
			return
		}
	}
}
