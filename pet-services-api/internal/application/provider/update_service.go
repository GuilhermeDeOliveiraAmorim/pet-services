package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// UpdateServiceUseCase atualiza um serviço existente do prestador.
type UpdateServiceUseCase struct {
	providerRepo provider.Repository
}

// NewUpdateServiceUseCase cria nova instância do caso de uso.
func NewUpdateServiceUseCase(providerRepo provider.Repository) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{providerRepo: providerRepo}
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
	if input.Name == "" {
		return provider.NewValidationError("service.name", "nome do serviço é obrigatório")
	}
	if input.PriceMin < 0 || input.PriceMax < 0 {
		return provider.NewValidationError("service.price", "preço não pode ser negativo")
	}
	if input.PriceMax > 0 && input.PriceMin > input.PriceMax {
		return provider.NewValidationError("service.price", "preço máximo deve ser maior ou igual ao mínimo")
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
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
		return provider.ErrServiceNotFound
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao salvar serviço: %w", err)
	}

	return nil
}
