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

func haversineDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	toRad := func(v float64) float64 {
		return v * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

type SearchServicesInput struct {
	Query      string
	CategoryID string
	TagID      string
	Latitude   float64
	Longitude  float64
	RadiusKm   float64
	PriceMin   float64
	PriceMax   float64
	Page       int
	PageSize   int
}

type SearchServicesOutput struct {
	Services     []*ServiceListItem `json:"services"`
	TotalItems   int64              `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	ItemsPerPage int                `json:"items_per_page"`
	SearchQuery  string             `json:"search_query,omitempty"`
	Location     *SearchLocation    `json:"location,omitempty"`
}

type SearchLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusKm  float64 `json:"radius_km"`
}

type SearchServicesUseCase struct {
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewSearchServicesUseCase(
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *SearchServicesUseCase {
	return &SearchServicesUseCase{
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *SearchServicesUseCase) Execute(ctx context.Context, input SearchServicesInput) (*SearchServicesOutput, []exceptions.ProblemDetails) {
	const from = "SearchServicesUseCase.Execute"

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	// Validação de parâmetros de localização
	if (input.Latitude != 0 || input.Longitude != 0) && input.RadiusKm <= 0 {
		input.RadiusKm = 10 // Default: 10km de raio
	}

	services, total, err := uc.serviceRepository.Search(
		input.Query,
		input.CategoryID,
		input.TagID,
		input.Latitude,
		input.Longitude,
		input.RadiusKm,
		input.PriceMin,
		input.PriceMax,
		input.Page,
		input.PageSize,
	)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviços", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	serviceItems := make([]*ServiceListItem, len(services))
	for i, svc := range services {
		// Buscar dados do provider
		businessName := ""
		averageRating := 0.0
		reviewCount := 0
		var distanceKm *float64
		provider, err := uc.providerRepository.FindByID(svc.ProviderID)
		if err == nil && provider != nil {
			businessName = provider.BusinessName
			averageRating = provider.AverageRating
			reviewCount = len(provider.Reviews)

			if input.Latitude != 0 && input.Longitude != 0 {
				providerLat := provider.Address.Location.Latitude
				providerLon := provider.Address.Location.Longitude
				if providerLat != 0 || providerLon != 0 {
					d := haversineDistanceKm(input.Latitude, input.Longitude, providerLat, providerLon)
					distanceKm = &d
				}
			}
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
			DistanceKm:    distanceKm,
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

	output := &SearchServicesOutput{
		Services:     serviceItems,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}

	if input.Query != "" {
		output.SearchQuery = input.Query
	}

	if input.Latitude != 0 && input.Longitude != 0 {
		output.Location = &SearchLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
			RadiusKm:  input.RadiusKm,
		}
	}

	return output, nil
}
