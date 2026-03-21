package usecases

import (
	"context"
	"errors"
	"math"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type GetProviderInput struct {
	ProviderID string `json:"provider_id"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type GetProviderOutput struct {
	Provider     *entities.Provider  `json:"provider"`
	Services     []*entities.Service `json:"services"`
	TotalItems   int64               `json:"total_items"`
	TotalPages   int                 `json:"total_pages"`
	CurrentPage  int                 `json:"current_page"`
	ItemsPerPage int                 `json:"items_per_page"`
}

type GetProviderUseCase struct {
	providerRepository entities.ProviderRepository
	serviceRepository  entities.ServiceRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewGetProviderUseCase(
	providerRepository entities.ProviderRepository,
	serviceRepository entities.ServiceRepository,
	storageService storage.ObjectStorage,
	logger logging.LoggerInterface,
) *GetProviderUseCase {
	return &GetProviderUseCase{
		providerRepository: providerRepository,
		serviceRepository:  serviceRepository,
		storage:            storageService,
		logger:             logger,
	}
}

func (uc *GetProviderUseCase) Execute(ctx context.Context, input GetProviderInput) (*GetProviderOutput, []exceptions.ProblemDetails) {
	const from = "GetProviderUseCase.Execute"

	if input.ProviderID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do provedor ausente", errors.New("O ID do provedor é obrigatório"))
	}

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	provider, err := uc.providerRepository.FindByID(input.ProviderID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	if !provider.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Provedor inativo", errors.New("O provedor informado está inativo"))
	}

	if len(provider.Photos) > 0 {
		for i := range provider.Photos {
			key := provider.Photos[i].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "providers/" + provider.ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					provider.Photos[i].URL = url
				}
			}
		}
	}

	services, total, err := uc.serviceRepository.List(provider.ID, "", "", 0, 0, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar serviços do provedor", err)
	}

	for i := range services {
		if len(services[i].Photos) == 0 {
			continue
		}
		for j := range services[i].Photos {
			key := services[i].Photos[j].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "services/" + services[i].ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					services[i].Photos[j].URL = url
				}
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	return &GetProviderOutput{
		Provider:     provider,
		Services:     services,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}, nil
}
