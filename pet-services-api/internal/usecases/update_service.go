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

type UpdateServiceInputBody struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	PriceMinimum *float64 `json:"price_minimum,omitempty"`
	PriceMaximum *float64 `json:"price_maximum,omitempty"`
	Duration     *int     `json:"duration,omitempty"`
}

type UpdateServiceInput struct {
	UserID       string   `json:"user_id"`
	ServiceID    string   `json:"service_id"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	PriceMinimum *float64 `json:"price_minimum,omitempty"`
	PriceMaximum *float64 `json:"price_maximum,omitempty"`
	Duration     *int     `json:"duration,omitempty"`
}

type UpdateServiceOutput struct {
	Message string            `json:"message,omitempty"`
	Detail  string            `json:"detail,omitempty"`
	Service *entities.Service `json:"service,omitempty"`
}

type UpdateServiceUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	serviceRepository  entities.ServiceRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewUpdateServiceUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	serviceRepository entities.ServiceRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		serviceRepository:  serviceRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *UpdateServiceUseCase) Execute(ctx context.Context, input UpdateServiceInput) (*UpdateServiceOutput, []exceptions.ProblemDetails) {
	const from = "UpdateServiceUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem atualizar serviços"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	service, err := uc.serviceRepository.FindByID(input.ServiceID)
	if err != nil {
		if err.Error() == consts.ServiceNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Serviço não encontrado", errors.New("Não foi possível encontrar o serviço informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviço", err)
	}

	if service.ProviderID != provider.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O serviço informado não pertence ao provedor autenticado"))
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço inativo", errors.New("O serviço informado está inativo"))
	}

	if input.Name != "" {
		if len(input.Name) > 100 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Nome muito longo", errors.New("O nome do serviço deve ter no máximo 100 caracteres"))
		}
		service.Name = input.Name
	}

	if input.Description != "" {
		if len(input.Description) > 1000 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Descrição muito longa", errors.New("A descrição deve ter no máximo 1000 caracteres"))
		}
		service.Description = input.Description
	}

	if input.Price != nil && *input.Price < 0 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Preço do serviço inválido", errors.New("O preço do serviço não pode ser negativo"))
	}

	if input.PriceMinimum != nil && *input.PriceMinimum < 0 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Preço mínimo do serviço inválido", errors.New("O preço mínimo do serviço não pode ser negativo"))
	}

	if input.PriceMaximum != nil && *input.PriceMaximum < 0 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Preço máximo do serviço inválido", errors.New("O preço máximo do serviço não pode ser negativo"))
	}

	if input.Price != nil {
		service.Price = *input.Price
	}

	if input.PriceMinimum != nil {
		service.PriceMinimum = *input.PriceMinimum
	}

	if input.PriceMaximum != nil {
		service.PriceMaximum = *input.PriceMaximum
	}

	if service.Price > 0 && (service.PriceMinimum > 0 || service.PriceMaximum > 0) {
		return nil, uc.logger.LogBadRequest(ctx, from, "Conflito de preços", errors.New("Use 'price' para preço fixo OU 'price_minimum' e 'price_maximum' para faixa de preço, não ambos"))
	}

	if service.PriceMinimum > 0 && service.PriceMaximum > 0 && service.PriceMinimum > service.PriceMaximum {
		return nil, uc.logger.LogBadRequest(ctx, from, "Faixa de preço inválida", errors.New("O preço mínimo não pode ser maior que o preço máximo"))
	}

	if service.Price == 0 && service.PriceMinimum == 0 && service.PriceMaximum == 0 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Preço ausente", errors.New("Defina um preço fixo ou uma faixa de preço para o serviço"))
	}

	if input.Duration != nil && *input.Duration < 0 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Duração do serviço inválida", errors.New("A duração do serviço não pode ser negativa"))
	}

	if input.Duration != nil {
		service.Duration = *input.Duration
	}

	service.Updated()

	if err := uc.serviceRepository.Update(service); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar serviço", err)
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

	return &UpdateServiceOutput{
		Message: "Serviço atualizado com sucesso",
		Detail:  "Os dados do serviço foram atualizados",
		Service: service,
	}, nil
}
