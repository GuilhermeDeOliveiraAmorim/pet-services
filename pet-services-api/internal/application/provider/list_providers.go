package provider

import (
	"context"
	"fmt"

	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ListProvidersByLocationUseCase busca prestadores ativos por localização/raio.
type ListProvidersByLocationUseCase struct {
	providerRepo provider.Repository
}

// NewListProvidersByLocationUseCase cria instância do caso de uso.
func NewListProvidersByLocationUseCase(providerRepo provider.Repository) *ListProvidersByLocationUseCase {
	return &ListProvidersByLocationUseCase{providerRepo: providerRepo}
}

// ListProvidersByLocationInput filtros de busca.
type ListProvidersByLocationInput struct {
	Latitude  float64
	Longitude float64
	RadiusKM  float64
	Page      int
	Limit     int
}

// Execute retorna prestadores ativos próximos com paginação.
func (uc *ListProvidersByLocationUseCase) Execute(ctx context.Context, input ListProvidersByLocationInput) ([]*provider.Provider, int64, error) {
	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return nil, 0, provider.ErrInvalidLocation
	}
	if input.RadiusKM <= 0 {
		return nil, 0, fmt.Errorf("raio deve ser maior que zero")
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	results, total, err := uc.providerRepo.FindActiveByLocation(ctx, input.Latitude, input.Longitude, input.RadiusKM, page, limit)
	if err != nil {
		return nil, 0, err
	}

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
