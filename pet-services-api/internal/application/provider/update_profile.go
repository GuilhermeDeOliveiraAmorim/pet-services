package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
	"github.com/guilherme/pet-services-api/internal/domain/user"
)

// UpdateProviderProfileUseCase atualiza dados básicos do prestador.
type UpdateProviderProfileUseCase struct {
	providerRepo provider.Repository
}

// NewUpdateProviderProfileUseCase cria nova instância.
func NewUpdateProviderProfileUseCase(providerRepo provider.Repository) *UpdateProviderProfileUseCase {
	return &UpdateProviderProfileUseCase{providerRepo: providerRepo}
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
}

// Execute aplica alterações parciais e persiste.
func (uc *UpdateProviderProfileUseCase) Execute(ctx context.Context, input UpdateProviderProfileInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if input.UserID != uuid.Nil && p.UserID != input.UserID {
		return fmt.Errorf("não autorizado a atualizar este perfil")
	}

	if input.BusinessName != nil {
		name := strings.TrimSpace(*input.BusinessName)
		if name == "" {
			return provider.NewValidationError("business_name", "nome do negócio é obrigatório")
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
			return provider.ErrInvalidPriceRange
		}
	}

	locationUpdate := input.Latitude != nil || input.Longitude != nil || input.Address != nil
	if locationUpdate {
		if input.Latitude == nil || input.Longitude == nil {
			return provider.ErrInvalidLocation
		}
		lat := *input.Latitude
		lon := *input.Longitude
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			return provider.ErrInvalidLocation
		}

		addr := p.Address
		if input.Address != nil {
			addr = *input.Address
		}

		p.SetLocation(lat, lon, addr)
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao atualizar perfil: %w", err)
	}

	return nil
}
