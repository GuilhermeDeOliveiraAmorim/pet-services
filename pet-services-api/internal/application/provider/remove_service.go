package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// RemoveServiceUseCase remove um serviço existente do prestador.
type RemoveServiceUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewRemoveServiceUseCase cria nova instância do caso de uso.
func NewRemoveServiceUseCase(providerRepo provider.Repository, logger *slog.Logger) *RemoveServiceUseCase {
	return &RemoveServiceUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// RemoveServiceInput representa os dados para remover serviço.
type RemoveServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
}

// Execute remove o serviço do prestador.
func (uc *RemoveServiceUseCase) Execute(ctx context.Context, input RemoveServiceInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "RemoveServiceUseCase", slog.String("provider_id", input.ProviderID.String()), slog.String("service", input.Name))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	before := len(p.Services)
	p.RemoveService(input.Category, input.Name)

	if len(p.Services) == before {
		err = provider.ErrServiceNotFound
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao remover serviço: %w", err)
		return err
	}

	return nil
}
