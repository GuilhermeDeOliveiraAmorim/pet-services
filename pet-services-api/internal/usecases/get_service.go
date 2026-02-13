package usecases

import (
	"context"
	"errors"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type GetServiceInput struct {
	ServiceID string `json:"service_id"`
}

type GetServiceOutput struct {
	Service *entities.Service `json:"service"`
}

type GetServiceUseCase struct {
	serviceRepository entities.ServiceRepository
	storage           storage.ObjectStorage
	logger            logging.LoggerInterface
}

func NewGetServiceUseCase(
	serviceRepository entities.ServiceRepository,
	storageService storage.ObjectStorage,
	logger logging.LoggerInterface,
) *GetServiceUseCase {
	return &GetServiceUseCase{
		serviceRepository: serviceRepository,
		storage:           storageService,
		logger:            logger,
	}
}

func (uc *GetServiceUseCase) Execute(ctx context.Context, input GetServiceInput) (*GetServiceOutput, []exceptions.ProblemDetails) {
	const from = "GetServiceUseCase.Execute"

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	service, err := uc.serviceRepository.FindByID(input.ServiceID)
	if err != nil {
		if err.Error() == consts.ServiceNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Serviço não encontrado", errors.New("Não foi possível encontrar o serviço informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviço", err)
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço inativo", errors.New("O serviço informado está inativo"))
	}

	if len(service.Photos) > 0 {
		for i := range service.Photos {
			key := service.Photos[i].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "services/" + service.ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					service.Photos[i].URL = url
				}
			}
		}
	}

	return &GetServiceOutput{
		Service: service,
	}, nil
}
