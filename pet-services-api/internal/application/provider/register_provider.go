package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
	"github.com/guilherme/pet-services-api/internal/domain/user"
)

// RegisterProviderUseCase orquestra a criação de um perfil de prestador.
type RegisterProviderUseCase struct {
	providerRepo provider.Repository
	userRepo     user.Repository
}

// NewRegisterProviderUseCase cria um novo caso de uso.
func NewRegisterProviderUseCase(providerRepo provider.Repository, userRepo user.Repository) *RegisterProviderUseCase {
	return &RegisterProviderUseCase{
		providerRepo: providerRepo,
		userRepo:     userRepo,
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
	if input.BusinessName == "" {
		return nil, provider.NewValidationError("business_name", "nome do negócio é obrigatório")
	}

	if len(input.Services) == 0 {
		return nil, provider.ErrNoServicesProvided
	}

	for _, svc := range input.Services {
		if svc.Name == "" {
			return nil, provider.NewValidationError("service.name", "nome do serviço é obrigatório")
		}
		if svc.PriceMin < 0 || svc.PriceMax < 0 {
			return nil, provider.NewValidationError("service.price", "preço não pode ser negativo")
		}
		if svc.PriceMax > 0 && svc.PriceMin > svc.PriceMax {
			return nil, provider.NewValidationError("service.price", "preço máximo deve ser maior ou igual ao mínimo")
		}
	}

	// Verifica se o usuário existe e é do tipo provider
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}
	if u.Type != user.UserTypeProvider {
		return nil, fmt.Errorf("usuário não é do tipo provider")
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

	return &RegisterProviderOutput{
		ProviderID:   p.ID,
		BusinessName: p.BusinessName,
		IsActive:     p.IsActive,
	}, nil
}
