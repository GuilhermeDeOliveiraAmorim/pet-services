package request

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// CreateRequestUseCase registra uma nova solicitação de serviço.
type CreateRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       *slog.Logger
}

// NewCreateRequestUseCase cria uma nova instância.
func NewCreateRequestUseCase(
	requestRepo domainRequest.Repository,
	providerRepo domainProvider.Repository,
	logger *slog.Logger,
) *CreateRequestUseCase {
	return &CreateRequestUseCase{
		requestRepo:  requestRepo,
		providerRepo: providerRepo,
		logger:       logging.EnsureLogger(logger),
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
	var (
		result *domainRequest.ServiceRequest
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "CreateRequestUseCase", slog.String("owner_id", input.OwnerID.String()), slog.String("provider_id", input.ProviderID.String()), slog.String("service_type", input.ServiceType))(&err)

	if input.OwnerID == uuid.Nil {
		err = fmt.Errorf("ownerID é obrigatório")
		return nil, err
	}
	if input.ProviderID == uuid.Nil {
		err = domainProvider.ErrProviderNotFound
		return nil, err
	}
	if strings.TrimSpace(input.ServiceType) == "" {
		err = domainRequest.ErrInvalidServiceType
		return nil, err
	}

	// Validar data preferencial: não pode ser no passado
	today := time.Now()
	if input.PreferredDate.IsZero() || input.PreferredDate.Before(today.Truncate(24*time.Hour)) {
		err = domainRequest.ErrInvalidPreferredDate
		return nil, err
	}

	// Validar horário
	if strings.TrimSpace(input.PreferredTime) == "" {
		err = fmt.Errorf("horário preferencial é obrigatório")
		return nil, err
	}
	if _, err := time.Parse("15:04", input.PreferredTime); err != nil {
		err = fmt.Errorf("horário preferencial inválido, use HH:MM")
		return nil, err
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
		err = domainProvider.ErrProviderNotFound
		return nil, err
	}
	if !provider.IsActive {
		err = domainProvider.ErrProviderNotActive
		return nil, err
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
		err = fmt.Errorf("falha ao salvar solicitação: %w", err)
		return nil, err
	}

	result = req

	return result, nil
}
