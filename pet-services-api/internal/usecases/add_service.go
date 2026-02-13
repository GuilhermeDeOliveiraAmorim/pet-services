package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type AddServiceInputBody struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	PriceMinimum float64 `json:"price_minimum"`
	PriceMaximum float64 `json:"price_maximum"`
	Duration     int     `json:"duration"`
}

type AddServiceInput struct {
	UserID       string  `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	PriceMinimum float64 `json:"price_minimum"`
	PriceMaximum float64 `json:"price_maximum"`
	Duration     int     `json:"duration"`
}

type AddServiceOutput struct {
	Message string           `json:"message,omitempty"`
	Detail  string           `json:"detail,omitempty"`
	Service *entities.Service `json:"service,omitempty"`
}

type AddServiceUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	serviceRepository  entities.ServiceRepository
	logger             logging.LoggerInterface
}

func NewAddServiceUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	serviceRepository entities.ServiceRepository,
	logger logging.LoggerInterface,
) *AddServiceUseCase {
	return &AddServiceUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		serviceRepository:  serviceRepository,
		logger:             logger,
	}
}

func (uc *AddServiceUseCase) Execute(ctx context.Context, input AddServiceInput) (*AddServiceOutput, []exceptions.ProblemDetails) {
	const from = "AddServiceUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem adicionar serviços"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	service, problems := entities.NewService(
		provider.ID,
		input.Name,
		input.Description,
		input.Price,
		input.PriceMinimum,
		input.PriceMaximum,
		input.Duration,
	)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Serviço inválido", problems)
		return nil, problems
	}

	if err := uc.serviceRepository.Create(service); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao adicionar serviço", err)
	}

	return &AddServiceOutput{
		Message: "Serviço adicionado com sucesso",
		Detail:  "O serviço foi cadastrado para o provedor",
		Service: service,
	}, nil
}
