package usecases

import (
	"context"
	"math"
	"strings"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type ListServicesInput struct {
	ProviderID string
	CategoryID string
	TagID      string
	PriceMin   float64
	PriceMax   float64
	Page       int
	PageSize   int
}

type ServiceListItem struct {
	ID            string              `json:"id"`
	ProviderID    string              `json:"provider_id"`
	BusinessName  string              `json:"business_name"`
	AverageRating float64             `json:"average_rating"`
	ReviewCount   int                 `json:"review_count"`
	DistanceKm    *float64            `json:"distance_km,omitempty"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Price         float64             `json:"price"`
	PriceMinimum  float64             `json:"price_minimum"`
	PriceMaximum  float64             `json:"price_maximum"`
	Duration      int                 `json:"duration"`
	Photos        []entities.Photo    `json:"photos"`
	Categories    []entities.Category `json:"categories"`
	Tags          []entities.Tag      `json:"tags"`
	Active        bool                `json:"active"`
	CreatedAt     string              `json:"created_at"`
}

type ListServicesOutput struct {
	Services     []*ServiceListItem `json:"services"`
	TotalItems   int64              `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	ItemsPerPage int                `json:"items_per_page"`
}

type ListServicesUseCase struct {
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewListServicesUseCase(
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *ListServicesUseCase {
	return &ListServicesUseCase{
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *ListServicesUseCase) Execute(ctx context.Context, input ListServicesInput) (*ListServicesOutput, []exceptions.ProblemDetails) {
	const from = "ListServicesUseCase.Execute"

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	services, total, err := uc.serviceRepository.List(
		input.ProviderID,
		input.CategoryID,
		input.TagID,
		input.PriceMin,
		input.PriceMax,
		input.Page,
		input.PageSize,
	)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar serviços", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	serviceItems := make([]*ServiceListItem, len(services))
	for i, svc := range services {
		// Buscar dados do provider
		businessName := ""
		averageRating := 0.0
		reviewCount := 0
		provider, err := uc.providerRepository.FindByID(svc.ProviderID)
		if err == nil && provider != nil {
			businessName = provider.BusinessName
			averageRating = provider.AverageRating
			reviewCount = len(provider.Reviews)
		}

		// Assinar URLs das fotos
		if len(svc.Photos) > 0 {
			for j := range svc.Photos {
				key := svc.Photos[j].URL
				if key != "" && !strings.HasPrefix(key, "http") {
					if !strings.Contains(key, "/") {
						key = "services/" + svc.ID + "/" + key
					}
					url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
					if err == nil {
						svc.Photos[j].URL = url
					}
				}
			}
		}

		serviceItems[i] = &ServiceListItem{
			ID:            svc.ID,
			ProviderID:    svc.ProviderID,
			BusinessName:  businessName,
			AverageRating: averageRating,
			ReviewCount:   reviewCount,
			Name:          svc.Name,
			Description:   svc.Description,
			Price:         svc.Price,
			PriceMinimum:  svc.PriceMinimum,
			PriceMaximum:  svc.PriceMaximum,
			Duration:      svc.Duration,
			Photos:        svc.Photos,
			Categories:    svc.Categories,
			Tags:          svc.Tags,
			Active:        svc.Active,
			CreatedAt:     svc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &ListServicesOutput{
		Services:     serviceItems,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}, nil
}
