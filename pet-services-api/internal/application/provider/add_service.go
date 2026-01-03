package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// AddServiceUseCase adiciona um serviço ao prestador.
type AddServiceUseCase struct {
	providerRepo provider.Repository
}

// NewAddServiceUseCase cria uma nova instância do caso de uso.
func NewAddServiceUseCase(providerRepo provider.Repository) *AddServiceUseCase {
	return &AddServiceUseCase{providerRepo: providerRepo}
}

// AddServiceInput representa os dados para adicionar serviço.
type AddServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
	PriceMin   float64
	PriceMax   float64
}

// Execute valida e adiciona o serviço ao prestador.
func (uc *AddServiceUseCase) Execute(ctx context.Context, input AddServiceInput) error {
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

	// Checagem simples para evitar duplicidade (categoria + nome)
	for _, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			return provider.NewValidationError("service", "serviço já cadastrado")
		}
	}

	p.AddService(input.Category, input.Name, input.PriceMin, input.PriceMax)

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao salvar serviço: %w", err)
	}

	return nil
}
