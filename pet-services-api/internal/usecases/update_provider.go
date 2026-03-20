package usecases

import (
	"context"
	"errors"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type UpdateProviderInputBody struct {
	BusinessName string            `json:"business_name,omitempty"`
	Description  string            `json:"description,omitempty"`
	PriceRange   string            `json:"price_range,omitempty"`
	Address      *entities.Address `json:"address,omitempty"`
}

type UpdateProviderInput struct {
	UserID       string            `json:"user_id"`
	ProviderID   string            `json:"provider_id"`
	BusinessName string            `json:"business_name,omitempty"`
	Description  string            `json:"description,omitempty"`
	PriceRange   string            `json:"price_range,omitempty"`
	Address      *entities.Address `json:"address,omitempty"`
}

type UpdateProviderOutput struct {
	Message  string             `json:"message,omitempty"`
	Detail   string             `json:"detail,omitempty"`
	Provider *entities.Provider `json:"provider,omitempty"`
}

type UpdateProviderUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	logger             logging.LoggerInterface
}

func NewUpdateProviderUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	logger logging.LoggerInterface,
) *UpdateProviderUseCase {
	return &UpdateProviderUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		logger:             logger,
	}
}

func (uc *UpdateProviderUseCase) Execute(ctx context.Context, input UpdateProviderInput) (*UpdateProviderOutput, []exceptions.ProblemDetails) {
	const from = "UpdateProviderUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ProviderID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do provedor ausente", errors.New("O ID do provedor é obrigatório"))
	}

	if strings.TrimSpace(input.BusinessName) == "" &&
		strings.TrimSpace(input.Description) == "" &&
		strings.TrimSpace(input.PriceRange) == "" &&
		input.Address == nil {
		return nil, uc.logger.LogBadRequest(ctx, from, "Nenhum campo para atualização", errors.New("Informe ao menos um campo para atualizar"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem atualizar provedor"))
	}

	provider, err := uc.providerRepository.FindByID(input.ProviderID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	if provider.UserID != user.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O provedor não pertence ao usuário autenticado"))
	}

	if businessName := strings.TrimSpace(input.BusinessName); businessName != "" {
		provider.BusinessName = businessName
	}

	if description := strings.TrimSpace(input.Description); description != "" {
		provider.Description = description
	}

	if priceRange := strings.TrimSpace(input.PriceRange); priceRange != "" {
		provider.PriceRange = priceRange
	}

	if input.Address != nil {
		addressInput := input.Address

		location, errs := entities.NewLocation(
			addressInput.Location.Latitude,
			addressInput.Location.Longitude,
		)
		if len(errs) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Localização inválida", errs)
			return nil, errs
		}

		address, errs := entities.NewAddress(
			addressInput.Street,
			addressInput.Number,
			addressInput.Neighborhood,
			addressInput.City,
			addressInput.ZipCode,
			addressInput.State,
			addressInput.Country,
			addressInput.Complement,
			*location,
		)
		if len(errs) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Endereço inválido", errs)
			return nil, errs
		}

		provider.Address = *address
	}

	provider.Updated()

	if err := uc.providerRepository.Update(provider); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar provedor", err)
	}

	return &UpdateProviderOutput{
		Message:  "Provedor atualizado com sucesso",
		Detail:   "Os dados gerais do provedor foram atualizados",
		Provider: provider,
	}, nil
}
