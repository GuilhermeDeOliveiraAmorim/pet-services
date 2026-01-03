package provider

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
	"github.com/guilherme/pet-services-api/internal/domain/user"
)

// UpdateProviderProfileUseCase atualiza dados básicos do prestador.
type UpdateProviderProfileUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewUpdateProviderProfileUseCase cria nova instância.
func NewUpdateProviderProfileUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateProviderProfileUseCase {
	return &UpdateProviderProfileUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
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
	var err error
	defer logging.UseCase(ctx, uc.logger, "UpdateProviderProfileUseCase", slog.String("provider_id", input.ProviderID.String()), slog.String("user_id", input.UserID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if input.UserID != uuid.Nil && p.UserID != input.UserID {
		err = fmt.Errorf("não autorizado a atualizar este perfil")
		return err
	}

	if input.BusinessName != nil {
		name := strings.TrimSpace(*input.BusinessName)
		if name == "" {
			err = provider.NewValidationError("business_name", "nome do negócio é obrigatório")
			return err
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
			err = provider.ErrInvalidPriceRange
			return err
		}
	}

	locationUpdate := input.Latitude != nil || input.Longitude != nil || input.Address != nil
	if locationUpdate {
		if input.Latitude == nil || input.Longitude == nil {
			err = provider.ErrInvalidLocation
			return err
		}
		lat := *input.Latitude
		lon := *input.Longitude
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			err = provider.ErrInvalidLocation
			return err
		}

		addr := p.Address
		if input.Address != nil {
			addr = *input.Address
		}

		p.SetLocation(lat, lon, addr)
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao atualizar perfil: %w", err)
		return err
	}

	return nil
}
