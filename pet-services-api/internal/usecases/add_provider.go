package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type AddProviderInputBody struct {
	BusinessName string           `json:"business_name"`
	Description  string           `json:"description"`
	PriceRange   string           `json:"price_range"`
	Address      entities.Address `json:"address"`
}

type AddProviderInput struct {
	UserID       string           `json:"user_id"`
	BusinessName string           `json:"business_name"`
	Description  string           `json:"description"`
	PriceRange   string           `json:"price_range"`
	Address      entities.Address `json:"address"`
}

type AddProviderOutput struct {
	Message  string             `json:"message,omitempty"`
	Detail   string             `json:"detail,omitempty"`
	Provider *entities.Provider `json:"provider,omitempty"`
}

type AddProviderUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	logger             logging.LoggerInterface
}

func NewAddProviderUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	logger logging.LoggerInterface,
) *AddProviderUseCase {
	return &AddProviderUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		logger:             logger,
	}
}

func (uc *AddProviderUseCase) Execute(ctx context.Context, input AddProviderInput) (*AddProviderOutput, []exceptions.ProblemDetails) {
	const from = "AddProviderUseCase.Execute"

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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem criar provedor"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err == nil && provider != nil {
		return nil, uc.logger.LogConflict(ctx, from, "Provedor já existe", errors.New("O usuário já possui um provedor cadastrado"))
	}
	if err != nil && err.Error() != consts.ProviderNotFoundError {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	var problems []exceptions.ProblemDetails

	location, errs := entities.NewLocation(input.Address.Location.Latitude, input.Address.Location.Longitude)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Localização inválida", errs)
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	address, errs := entities.NewAddress(
		input.Address.Street,
		input.Address.Number,
		input.Address.Neighborhood,
		input.Address.City,
		input.Address.ZipCode,
		input.Address.State,
		input.Address.Country,
		input.Address.Complement,
		*location,
	)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Endereço inválido", errs)
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	providerEntity, errs := entities.NewProvider(
		user.ID,
		input.BusinessName,
		*address,
		input.Description,
		input.PriceRange,
	)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Provedor inválido", errs)
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	if err := uc.providerRepository.Create(providerEntity); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar provedor", err)
	}

	return &AddProviderOutput{
		Message:  "Provedor criado com sucesso",
		Detail:   "O provedor foi cadastrado para o usuário",
		Provider: providerEntity,
	}, nil
}
