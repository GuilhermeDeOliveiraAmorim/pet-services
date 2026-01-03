package request

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// CreateRequestUseCase registra uma nova solicitação de serviço.
type CreateRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
}

// NewCreateRequestUseCase cria uma nova instância.
func NewCreateRequestUseCase(
	requestRepo domainRequest.Repository,
	providerRepo domainProvider.Repository,
) *CreateRequestUseCase {
	return &CreateRequestUseCase{
		requestRepo:  requestRepo,
		providerRepo: providerRepo,
	}
}

// CreateRequestInput representa os dados necessários para criar uma solicitação.
type CreateRequestInput struct {
	OwnerID       uuid.UUID
	ProviderID    uuid.UUID
	ServiceType   string
	PetName       string
	PetType       domainRequest.PetType
	PetBreed      string
	PetAge        int
	PetWeight     float64
	PetNotes      string
	PreferredDate time.Time
	PreferredTime string // HH:MM
	Notes         string
}

// Execute valida os dados, cria a entidade e persiste.
func (uc *CreateRequestUseCase) Execute(ctx context.Context, input CreateRequestInput) (*domainRequest.ServiceRequest, error) {
	if input.OwnerID == uuid.Nil {
		return nil, fmt.Errorf("ownerID é obrigatório")
	}
	if input.ProviderID == uuid.Nil {
		return nil, domainProvider.ErrProviderNotFound
	}
	if strings.TrimSpace(input.ServiceType) == "" {
		return nil, domainRequest.ErrInvalidServiceType
	}

	// Validar data preferencial: não pode ser no passado
	today := time.Now()
	if input.PreferredDate.IsZero() || input.PreferredDate.Before(today.Truncate(24*time.Hour)) {
		return nil, domainRequest.ErrInvalidPreferredDate
	}

	// Validar horário
	if strings.TrimSpace(input.PreferredTime) == "" {
		return nil, fmt.Errorf("horário preferencial é obrigatório")
	}
	if _, err := time.Parse("15:04", input.PreferredTime); err != nil {
		return nil, fmt.Errorf("horário preferencial inválido, use HH:MM")
	}

	petInfo, err := domainRequest.NewPetInfo(
		input.PetName,
		input.PetType,
		input.PetBreed,
		input.PetAge,
		input.PetWeight,
		input.PetNotes,
	)
	if err != nil {
		return nil, err
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return nil, domainProvider.ErrProviderNotFound
	}
	if !provider.IsActive {
		return nil, domainProvider.ErrProviderNotActive
	}

	req := domainRequest.NewServiceRequest(
		input.OwnerID,
		input.ProviderID,
		strings.TrimSpace(input.ServiceType),
		petInfo,
		input.PreferredDate,
		input.PreferredTime,
		strings.TrimSpace(input.Notes),
	)

	if err := uc.requestRepo.Create(ctx, req); err != nil {
		return nil, fmt.Errorf("falha ao salvar solicitação: %w", err)
	}

	return req, nil
}
