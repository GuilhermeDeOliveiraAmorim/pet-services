package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// AddServiceUseCase adiciona um serviço ao prestador.
type AddServiceUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewAddServiceUseCase cria uma nova instância do caso de uso.
func NewAddServiceUseCase(providerRepo provider.Repository, logger *slog.Logger) *AddServiceUseCase {
	return &AddServiceUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
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
	var err error
	defer logging.UseCase(ctx, uc.logger, "AddServiceUseCase", slog.String("provider_id", input.ProviderID.String()), slog.String("service", input.Name))(&err)

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

	// Checagem simples para evitar duplicidade (categoria + nome)
	for _, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			err = provider.NewValidationError("service", "serviço já cadastrado")
			return err
		}
	}

	p.AddService(input.Category, input.Name, input.PriceMin, input.PriceMax)

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao salvar serviço: %w", err)
		return err
	}

	return nil
}
