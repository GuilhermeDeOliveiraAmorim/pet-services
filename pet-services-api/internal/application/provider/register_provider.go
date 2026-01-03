package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
	"github.com/guilherme/pet-services-api/internal/domain/user"
)

// RegisterProviderUseCase orquestra a criação de um perfil de prestador.
type RegisterProviderUseCase struct {
	providerRepo provider.Repository
	userRepo     user.Repository
	logger       *slog.Logger
}

// NewRegisterProviderUseCase cria um novo caso de uso.
func NewRegisterProviderUseCase(providerRepo provider.Repository, userRepo user.Repository, logger *slog.Logger) *RegisterProviderUseCase {
	return &RegisterProviderUseCase{
		providerRepo: providerRepo,
		userRepo:     userRepo,
		logger:       logging.EnsureLogger(logger),
	}
}

// RegisterProviderInput representa os dados necessários para criar um prestador.
type RegisterProviderInput struct {
	UserID       uuid.UUID
	BusinessName string
	Description  string
	Address      user.Address
	Latitude     float64
	Longitude    float64
	Services     []ServiceInput
	PriceRange   provider.PriceRange
}

// ServiceInput representa um serviço a ser cadastrado.
type ServiceInput struct {
	Category string
	Name     string
	PriceMin float64
	PriceMax float64
}

// RegisterProviderOutput representa o resultado da criação.
type RegisterProviderOutput struct {
	ProviderID   uuid.UUID
	BusinessName string
	IsActive     bool
}

// Execute cria um perfil de prestador validando regras básicas de domínio.
func (uc *RegisterProviderUseCase) Execute(ctx context.Context, input RegisterProviderInput) (*RegisterProviderOutput, error) {
	var (
		result *RegisterProviderOutput
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "RegisterProviderUseCase", slog.String("user_id", input.UserID.String()), slog.String("business_name", input.BusinessName))(&err)

	if input.BusinessName == "" {
		err = provider.NewValidationError("business_name", "nome do negócio é obrigatório")
		return nil, err
	}

	if len(input.Services) == 0 {
		err = provider.ErrNoServicesProvided
		return nil, err
	}

	for _, svc := range input.Services {
		if svc.Name == "" {
			err = provider.NewValidationError("service.name", "nome do serviço é obrigatório")
			return nil, err
		}
		if svc.PriceMin < 0 || svc.PriceMax < 0 {
			err = provider.NewValidationError("service.price", "preço não pode ser negativo")
			return nil, err
		}
		if svc.PriceMax > 0 && svc.PriceMin > svc.PriceMax {
			err = provider.NewValidationError("service.price", "preço máximo deve ser maior ou igual ao mínimo")
			return nil, err
		}
	}

	// Verifica se o usuário existe e é do tipo provider
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		err = user.ErrUserNotFound
		return nil, err
	}
	if u.Type != user.UserTypeProvider {
		err = fmt.Errorf("usuário não é do tipo provider")
		return nil, err
	}

	// Cria entidade Provider
	p := provider.NewProvider(input.UserID, input.BusinessName, input.Description)

	// Define localização
	p.SetLocation(input.Latitude, input.Longitude, input.Address)

	// Define serviços
	for _, svc := range input.Services {
		p.AddService(svc.Category, svc.Name, svc.PriceMin, svc.PriceMax)
	}

	// Define faixa de preço (opcional)
	if input.PriceRange != "" {
		p.PriceRange = input.PriceRange
	}

	// Persiste
	if err := uc.providerRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	result = &RegisterProviderOutput{
		ProviderID:   p.ID,
		BusinessName: p.BusinessName,
		IsActive:     p.IsActive,
	}

	return result, nil
}
