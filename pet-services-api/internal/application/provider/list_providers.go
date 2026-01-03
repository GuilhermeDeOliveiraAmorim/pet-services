package provider

import (
	"context"
	"fmt"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ListProvidersByLocationUseCase busca prestadores ativos por localização/raio.
type ListProvidersByLocationUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewListProvidersByLocationUseCase cria instância do caso de uso.
func NewListProvidersByLocationUseCase(providerRepo provider.Repository, logger logging.LoggerService) *ListProvidersByLocationUseCase {
	return &ListProvidersByLocationUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// ListProvidersByLocationInput filtros de busca.
type ListProvidersByLocationInput struct {
	Latitude  float64
	Longitude float64
	RadiusKM  float64
	Page      int
	Limit     int
}

// Execute retorna prestadores ativos próximos com paginação seguindo padrão CreateRequestUseCase.
func (uc *ListProvidersByLocationUseCase) Execute(ctx context.Context, input ListProvidersByLocationInput) ([]*provider.Provider, int64, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ListProvidersByLocationUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "ListProvidersByLocationUseCase",
			Message: "Localização inválida",
			Error:   provider.ErrInvalidLocation,
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Localização inválida",
			Status: exceptions.RFC400_CODE,
			Detail: "Latitude ou longitude inválida.",
		}}
	}
	if input.RadiusKM <= 0 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "ListProvidersByLocationUseCase",
			Message: "Raio deve ser maior que zero",
			Error:   fmt.Errorf("raio deve ser maior que zero"),
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Raio inválido",
			Status: exceptions.RFC400_CODE,
			Detail: "O raio deve ser maior que zero.",
		}}
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	results, total, err := uc.providerRepo.FindActiveByLocation(ctx, input.Latitude, input.Longitude, input.RadiusKM, page, limit)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "ListProvidersByLocationUseCase",
			Message: "Falha ao buscar prestadores",
			Error:   err,
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao buscar prestadores",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ListProvidersByLocationUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return results, total, nil
}

// normalizePagination aplica defaults e evita valores inválidos.
func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
