package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// RemoveServiceUseCase remove um serviço existente do prestador.
type RemoveServiceUseCase struct {
	providerRepo provider.Repository
}

// NewRemoveServiceUseCase cria nova instância do caso de uso.
func NewRemoveServiceUseCase(providerRepo provider.Repository) *RemoveServiceUseCase {
	return &RemoveServiceUseCase{providerRepo: providerRepo}
}

// RemoveServiceInput representa os dados para remover serviço.
type RemoveServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
}

// Execute remove o serviço do prestador.
func (uc *RemoveServiceUseCase) Execute(ctx context.Context, input RemoveServiceInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	before := len(p.Services)
	p.RemoveService(input.Category, input.Name)

	if len(p.Services) == before {
		return provider.ErrServiceNotFound
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao remover serviço: %w", err)
	}

	return nil
}
