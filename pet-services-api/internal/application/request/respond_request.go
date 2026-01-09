package request

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainProvider "pet-services-api/internal/domain/provider"
	domainRequest "pet-services-api/internal/domain/request"
)

// AcceptRequestUseCase permite ao prestador aceitar uma solicitação.
type AcceptRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       logging.LoggerService
}

func NewAcceptRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger logging.LoggerService) *AcceptRequestUseCase {
	return &AcceptRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logger}
}

// AcceptRequestInput dados para aceitar uma solicitação.
type AcceptRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
}

const ACCEPT_REQUEST_USECASE = "ACCEPT_REQUEST_USECASE"

func (uc *AcceptRequestUseCase) Execute(ctx context.Context, input AcceptRequestInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    ACCEPT_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}
	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if req.ProviderID != input.ProviderID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Provider não autorizado para esta solicitação",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Provider não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "O prestador não está autorizado para esta solicitação.",
		}}
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
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
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Provider não está ativo",
			Error:   domainProvider.ErrProviderNotActive,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Provider não está ativo",
			Status: exceptions.RFC409_CODE,
			Detail: "O prestador informado não está ativo.",
		}}
	}

	if err := req.Accept(); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Solicitação não pode ser aceita",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Solicitação não pode ser aceita",
			Status: exceptions.RFC409_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    ACCEPT_REQUEST_USECASE,
			Message: "Falha ao aceitar solicitação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao aceitar solicitação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    ACCEPT_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}

// RejectRequestUseCase permite ao prestador rejeitar uma solicitação.
type RejectRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       logging.LoggerService
}

func NewRejectRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger logging.LoggerService) *RejectRequestUseCase {
	return &RejectRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logger}
}

// RejectRequestInput dados para rejeitar uma solicitação.
type RejectRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
	Reason     string
}

const REJECT_REQUEST_USECASE = "REJECT_REQUEST_USECASE"

func (uc *RejectRequestUseCase) Execute(ctx context.Context, input RejectRequestInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REJECT_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}
	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if req.ProviderID != input.ProviderID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Provider não autorizado para esta solicitação",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Provider não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "O prestador não está autorizado para esta solicitação.",
		}}
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
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
			From:    REJECT_REQUEST_USECASE,
			Message: "Provider não está ativo",
			Error:   domainProvider.ErrProviderNotActive,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Provider não está ativo",
			Status: exceptions.RFC409_CODE,
			Detail: "O prestador informado não está ativo.",
		}}
	}

	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Motivo da rejeição é obrigatório",
			Error:   errors.New("motivo da rejeição é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Motivo da rejeição é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O motivo da rejeição é obrigatório.",
		}}
	}

	if err := req.Reject(reason); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Solicitação não pode ser rejeitada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Solicitação não pode ser rejeitada",
			Status: exceptions.RFC409_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REJECT_REQUEST_USECASE,
			Message: "Falha ao rejeitar solicitação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao rejeitar solicitação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REJECT_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}

// CompleteRequestUseCase permite ao prestador concluir uma solicitação aceita.
type CompleteRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       logging.LoggerService
}

func NewCompleteRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger logging.LoggerService) *CompleteRequestUseCase {
	return &CompleteRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logger}
}

// CompleteRequestInput dados para concluir uma solicitação.
type CompleteRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
}

const COMPLETE_REQUEST_USECASE = "COMPLETE_REQUEST_USECASE"

func (uc *CompleteRequestUseCase) Execute(ctx context.Context, input CompleteRequestInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    COMPLETE_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}
	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if req.ProviderID != input.ProviderID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Provider não autorizado para esta solicitação",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Provider não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "O prestador não está autorizado para esta solicitação.",
		}}
	}

	if _, err := uc.providerRepo.FindByID(ctx, input.ProviderID); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Provider não encontrado",
			Error:   domainProvider.ErrProviderNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Provider não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	if err := req.Complete(); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Solicitação não pode ser concluída",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Solicitação não pode ser concluída",
			Status: exceptions.RFC409_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    COMPLETE_REQUEST_USECASE,
			Message: "Falha ao concluir solicitação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao concluir solicitação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    COMPLETE_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}

// CancelRequestUseCase permite ao dono cancelar uma solicitação.
type CancelRequestUseCase struct {
	requestRepo domainRequest.Repository
	logger      logging.LoggerService
}

func NewCancelRequestUseCase(requestRepo domainRequest.Repository, logger logging.LoggerService) *CancelRequestUseCase {
	return &CancelRequestUseCase{requestRepo: requestRepo, logger: logger}
}

// CancelRequestInput dados para cancelar uma solicitação.
type CancelRequestInput struct {
	RequestID uuid.UUID
	OwnerID   uuid.UUID
}

const CANCEL_REQUEST_USECASE = "CANCEL_REQUEST_USECASE"

func (uc *CancelRequestUseCase) Execute(ctx context.Context, input CancelRequestInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CANCEL_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}
	if input.OwnerID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "OwnerID é obrigatório",
			Error:   errors.New("ownerID é obrigatório"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "OwnerID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do dono é obrigatório.",
		}}
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if req.OwnerID != input.OwnerID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "Não autorizado a cancelar esta solicitação",
			Error:   errors.New("não autorizado a cancelar esta solicitação"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "Você não está autorizado a cancelar esta solicitação.",
		}}
	}

	if err := req.Cancel(); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "Solicitação não pode ser cancelada",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Solicitação não pode ser cancelada",
			Status: exceptions.RFC409_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CANCEL_REQUEST_USECASE,
			Message: "Falha ao cancelar solicitação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao cancelar solicitação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CANCEL_REQUEST_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
