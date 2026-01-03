package request

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// CreateRequestUseCase registra uma nova solicitação de serviço.
type CreateRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       logging.LoggerService
}

// NewCreateRequestUseCase cria uma nova instância.
func NewCreateRequestUseCase(
	requestRepo domainRequest.Repository,
	providerRepo domainProvider.Repository,
	logger logging.LoggerService,
) *CreateRequestUseCase {
	return &CreateRequestUseCase{
		requestRepo:  requestRepo,
		providerRepo: providerRepo,
		logger:       logger,
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

const CREATE_REQUEST_USECASE = "CREATE_REQUEST_USECASE"

// Execute valida os dados, cria a entidade e persiste, seguindo padrão de erros e logging.
func (uc *CreateRequestUseCase) Execute(ctx context.Context, input CreateRequestInput) (*domainRequest.ServiceRequest, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CREATE_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.OwnerID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "OwnerID é obrigatório",
			Error:   errors.New("ownerID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "OwnerID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do dono é obrigatório.",
		}}
	}
	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}
	if strings.TrimSpace(input.ServiceType) == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Tipo de serviço inválido",
			Error:   domainRequest.ErrInvalidServiceType,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Tipo de serviço inválido",
			Status: exceptions.RFC400_CODE,
			Detail: "O tipo de serviço é obrigatório.",
		}}
	}

	today := time.Now()
	if input.PreferredDate.IsZero() || input.PreferredDate.Before(today.Truncate(24*time.Hour)) {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Data preferencial inválida",
			Error:   domainRequest.ErrInvalidPreferredDate,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Data preferencial inválida",
			Status: exceptions.RFC400_CODE,
			Detail: "A data preferencial não pode ser no passado.",
		}}
	}

	if strings.TrimSpace(input.PreferredTime) == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Horário preferencial é obrigatório",
			Error:   errors.New("horário preferencial é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Horário preferencial é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O horário preferencial é obrigatório.",
		}}
	}
	if _, err := time.Parse("15:04", input.PreferredTime); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Horário preferencial inválido",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Horário preferencial inválido",
			Status: exceptions.RFC400_CODE,
			Detail: "O horário deve estar no formato HH:MM.",
		}}
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
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Dados do pet inválidos",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Dados do pet inválidos",
			Status: exceptions.RFC400_CODE,
			Detail: err.Error(),
		}}
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}
	if !provider.IsActive {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Provider não está ativo",
			Error:   domainProvider.ErrProviderNotActive,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Provider não está ativo",
			Status: exceptions.RFC409_CODE,
			Detail: "O prestador informado não está ativo.",
		}}
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
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CREATE_REQUEST_USECASE,
			Message: "Falha ao salvar solicitação",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar solicitação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CREATE_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return req, nil
}
