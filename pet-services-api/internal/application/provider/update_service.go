package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// UpdateServiceUseCase atualiza um serviço existente do prestador.
type UpdateServiceUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewUpdateServiceUseCase cria nova instância do caso de uso.
func NewUpdateServiceUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// UpdateServiceInput representa os dados para atualizar serviço.
type UpdateServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
	PriceMin   float64
	PriceMax   float64
}

// Execute valida e atualiza o serviço.
func (uc *UpdateServiceUseCase) Execute(ctx context.Context, input UpdateServiceInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "UpdateServiceUseCase", slog.String("provider_id", input.ProviderID.String()), slog.String("service", input.Name))(&err)

	if input.Name == "" {
		err = provider.NewValidationError("service.name", "nome do serviço é obrigatório")
		return err
	}
	if input.PriceMin < 0 || input.PriceMax < 0 {
		err = provider.NewValidationError("service.price", "preço não pode ser negativo")
		return err
	}
	if input.PriceMax > 0 && input.PriceMin > input.PriceMax {
		err = provider.NewValidationError("service.price", "preço máximo deve ser maior ou igual ao mínimo")
		return err
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	updated := false
	for i, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			p.Services[i].PriceMin = input.PriceMin
			p.Services[i].PriceMax = input.PriceMax
			updated = true
			break
		}
	}

	if !updated {
		err = provider.ErrServiceNotFound
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao salvar serviço: %w", err)
		return err
	}

	return nil
}
