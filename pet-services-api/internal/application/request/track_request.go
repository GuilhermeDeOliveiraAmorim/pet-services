package request

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainRequest "pet-services-api/internal/domain/request"
)

// GetRequestStatusUseCase retorna uma solicitação garantindo vínculo com dono ou prestador.

type GetRequestStatusUseCase struct {
	requestRepo domainRequest.Repository
	logger      logging.LoggerService
}

func NewGetRequestStatusUseCase(requestRepo domainRequest.Repository, logger logging.LoggerService) *GetRequestStatusUseCase {
	return &GetRequestStatusUseCase{requestRepo: requestRepo, logger: logger}
}

// GetRequestStatusInput entrada para buscar status.
type GetRequestStatusInput struct {
	RequestID  uuid.UUID
	OwnerID    uuid.UUID
	ProviderID uuid.UUID
}

const GET_REQUEST_STATUS_USECASE = "GET_REQUEST_STATUS_USECASE"

func (uc *GetRequestStatusUseCase) Execute(ctx context.Context, input GetRequestStatusInput) (*domainRequest.ServiceRequest, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    GET_REQUEST_STATUS_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    GET_REQUEST_STATUS_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}

	result, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    GET_REQUEST_STATUS_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if input.OwnerID != uuid.Nil && result.OwnerID == input.OwnerID {
		return result, nil
	}
	if input.ProviderID != uuid.Nil && result.ProviderID == input.ProviderID {
		return result, nil
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.ERROR,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC403_CODE,
		From:    GET_REQUEST_STATUS_USECASE,
		Message: "Não autorizado a visualizar esta solicitação",
		Error:   errors.New("não autorizado a visualizar esta solicitação"),
	})
	return nil, []exceptions.ProblemDetails{{
		Type:   exceptions.RFC403,
		Title:  "Não autorizado",
		Status: exceptions.RFC403_CODE,
		Detail: "Você não tem permissão para visualizar esta solicitação.",
	}}
}

// ListRequestsForOwnerUseCase lista solicitações do dono.
type ListRequestsForOwnerUseCase struct {
	requestRepo domainRequest.Repository
	logger      logging.LoggerService
}

func NewListRequestsForOwnerUseCase(requestRepo domainRequest.Repository, logger logging.LoggerService) *ListRequestsForOwnerUseCase {
	return &ListRequestsForOwnerUseCase{requestRepo: requestRepo, logger: logger}
}

// ListRequestsForOwnerInput filtro de listagem por dono.
type ListRequestsForOwnerInput struct {
	OwnerID uuid.UUID
	Page    int
	Limit   int
}

const LIST_REQUESTS_FOR_OWNER_USECASE = "LIST_REQUESTS_FOR_OWNER_USECASE"

func (uc *ListRequestsForOwnerUseCase) Execute(ctx context.Context, input ListRequestsForOwnerInput) ([]*domainRequest.ServiceRequest, int64, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_FOR_OWNER_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.OwnerID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    LIST_REQUESTS_FOR_OWNER_USECASE,
			Message: "OwnerID é obrigatório",
			Error:   errors.New("ownerID é obrigatório"),
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "OwnerID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do dono é obrigatório.",
		}}
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByOwnerID(ctx, input.OwnerID, page, limit)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LIST_REQUESTS_FOR_OWNER_USECASE,
			Message: "Erro ao buscar solicitações do dono",
			Error:   err,
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao buscar solicitações do dono",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_FOR_OWNER_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return requests, total, nil
}

// ListRequestsForProviderUseCase lista solicitações do prestador.
type ListRequestsForProviderUseCase struct {
	requestRepo domainRequest.Repository
	logger      logging.LoggerService
}

func NewListRequestsForProviderUseCase(requestRepo domainRequest.Repository, logger logging.LoggerService) *ListRequestsForProviderUseCase {
	return &ListRequestsForProviderUseCase{requestRepo: requestRepo, logger: logger}
}

// ListRequestsForProviderInput filtro de listagem por prestador.
type ListRequestsForProviderInput struct {
	ProviderID uuid.UUID
	Page       int
	Limit      int
}

const LIST_REQUESTS_FOR_PROVIDER_USECASE = "LIST_REQUESTS_FOR_PROVIDER_USECASE"

func (uc *ListRequestsForProviderUseCase) Execute(ctx context.Context, input ListRequestsForProviderInput) ([]*domainRequest.ServiceRequest, int64, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_FOR_PROVIDER_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    LIST_REQUESTS_FOR_PROVIDER_USECASE,
			Message: "ProviderID é obrigatório",
			Error:   errors.New("providerID é obrigatório"),
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "ProviderID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do prestador é obrigatório.",
		}}
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByProviderID(ctx, input.ProviderID, page, limit)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LIST_REQUESTS_FOR_PROVIDER_USECASE,
			Message: "Erro ao buscar solicitações do prestador",
			Error:   err,
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao buscar solicitações do prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_FOR_PROVIDER_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return requests, total, nil
}

// ListRequestsByStatusUseCase lista solicitações por status (uso administrativo ou dashboards).
type ListRequestsByStatusUseCase struct {
	requestRepo domainRequest.Repository
	logger      logging.LoggerService
}

func NewListRequestsByStatusUseCase(requestRepo domainRequest.Repository, logger logging.LoggerService) *ListRequestsByStatusUseCase {
	return &ListRequestsByStatusUseCase{requestRepo: requestRepo, logger: logger}
}

// ListRequestsByStatusInput filtro de listagem por status.
type ListRequestsByStatusInput struct {
	Status domainRequest.Status
	Page   int
	Limit  int
}

const LIST_REQUESTS_BY_STATUS_USECASE = "LIST_REQUESTS_BY_STATUS_USECASE"

func (uc *ListRequestsByStatusUseCase) Execute(ctx context.Context, input ListRequestsByStatusInput) ([]*domainRequest.ServiceRequest, int64, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_BY_STATUS_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.Status == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    LIST_REQUESTS_BY_STATUS_USECASE,
			Message: "Status é obrigatório",
			Error:   errors.New("status é obrigatório"),
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Status é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O status é obrigatório.",
		}}
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByStatus(ctx, input.Status, page, limit)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LIST_REQUESTS_BY_STATUS_USECASE,
			Message: "Erro ao buscar solicitações por status",
			Error:   err,
		})
		return nil, 0, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao buscar solicitações por status",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LIST_REQUESTS_BY_STATUS_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return requests, total, nil
}

// normalizePagination aplica defaults e evita valores inválidos.
func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
