package provider

import (
	"context"
	"fmt"
	"strings"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	"pet-services-api/internal/domain/provider"
	"pet-services-api/internal/domain/user"

	"github.com/google/uuid"
)

// UpdateProviderProfileUseCase atualiza dados básicos do prestador.
type UpdateProviderProfileUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewUpdateProviderProfileUseCase cria nova instância.
func NewUpdateProviderProfileUseCase(providerRepo provider.Repository, logger logging.LoggerService) *UpdateProviderProfileUseCase {
	return &UpdateProviderProfileUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// UpdateProviderProfileInput campos opcionais para atualização.
// Campos ponteiro indicam se devem ser alterados.
type UpdateProviderProfileInput struct {
	ProviderID   uuid.UUID
	UserID       uuid.UUID // para autorização simples (dono do perfil)
	BusinessName *string
	Description  *string
	Address      *user.Address
	Latitude     *float64
	Longitude    *float64
	PriceRange   *provider.PriceRange
	Photos       []string // URLs das novas fotos
}

// Execute aplica alterações parciais e persiste seguindo padrão CreateRequestUseCase.
func (uc *UpdateProviderProfileUseCase) Execute(ctx context.Context, input UpdateProviderProfileInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateProviderProfileUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateProviderProfileUseCase",
			Message: "ProviderID é obrigatório",
			Error:   fmt.Errorf("providerID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "ProviderID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do prestador é obrigatório.",
		}}
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "UpdateProviderProfileUseCase",
			Message: "Prestador não encontrado",
			Error:   provider.ErrProviderNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Prestador não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	if input.UserID != uuid.Nil && p.UserID != input.UserID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    "UpdateProviderProfileUseCase",
			Message: "Não autorizado a atualizar este perfil",
			Error:   fmt.Errorf("não autorizado a atualizar este perfil"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "Você não tem permissão para atualizar este perfil.",
		}}
	}

	if input.BusinessName != nil {
		name := strings.TrimSpace(*input.BusinessName)
		if name == "" {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "UpdateProviderProfileUseCase",
				Message: "Nome do negócio é obrigatório",
				Error:   provider.NewValidationError("business_name", "nome do negócio é obrigatório"),
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Nome do negócio é obrigatório",
				Status: exceptions.RFC400_CODE,
				Detail: "O nome do negócio é obrigatório.",
			}}
		}
		p.BusinessName = name
	}

	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		p.Description = desc
	}

	if input.PriceRange != nil {
		switch *input.PriceRange {
		case provider.PriceRangeLow, provider.PriceRangeMedium, provider.PriceRangeHigh:
			p.PriceRange = *input.PriceRange
		default:
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "UpdateProviderProfileUseCase",
				Message: "Faixa de preço inválida",
				Error:   provider.ErrInvalidPriceRange,
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Faixa de preço inválida",
				Status: exceptions.RFC400_CODE,
				Detail: "A faixa de preço informada é inválida.",
			}}
		}
	}

	locationUpdate := input.Latitude != nil || input.Longitude != nil || input.Address != nil
	if locationUpdate {
		if input.Latitude == nil || input.Longitude == nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "UpdateProviderProfileUseCase",
				Message: "Latitude e longitude são obrigatórios para atualizar localização",
				Error:   provider.ErrInvalidLocation,
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Localização inválida",
				Status: exceptions.RFC400_CODE,
				Detail: "Latitude e longitude são obrigatórios para atualizar localização.",
			}}
		}
		lat := *input.Latitude
		lon := *input.Longitude
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "UpdateProviderProfileUseCase",
				Message: "Latitude ou longitude inválida",
				Error:   provider.ErrInvalidLocation,
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Localização inválida",
				Status: exceptions.RFC400_CODE,
				Detail: "Latitude ou longitude inválida.",
			}}
		}

		addr := p.Address
		if input.Address != nil {
			addr = *input.Address
		}

		p.SetLocation(lat, lon, addr)
	}

	// Adiciona novas fotos sem remover as antigas
	for _, url := range input.Photos {
		_ = p.AddPhoto(url)
	}
	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "UpdateProviderProfileUseCase",
			Message: "Falha ao atualizar perfil",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar perfil",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateProviderProfileUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}
