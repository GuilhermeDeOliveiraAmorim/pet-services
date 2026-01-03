package request

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// GetRequestStatusUseCase retorna uma solicitação garantindo vínculo com dono ou prestador.
type GetRequestStatusUseCase struct {
	requestRepo domainRequest.Repository
	logger      *slog.Logger
}

func NewGetRequestStatusUseCase(requestRepo domainRequest.Repository, logger *slog.Logger) *GetRequestStatusUseCase {
	return &GetRequestStatusUseCase{requestRepo: requestRepo, logger: logging.EnsureLogger(logger)}
}

// GetRequestStatusInput entrada para buscar status.
type GetRequestStatusInput struct {
	RequestID  uuid.UUID
	OwnerID    uuid.UUID
	ProviderID uuid.UUID
}

func (uc *GetRequestStatusUseCase) Execute(ctx context.Context, input GetRequestStatusInput) (*domainRequest.ServiceRequest, error) {
	var (
		result *domainRequest.ServiceRequest
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "GetRequestStatusUseCase", slog.String("request_id", input.RequestID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return nil, err
	}

	result, err = uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	// Autorização simples: precisa ser dono ou prestador vinculado.
	if input.OwnerID != uuid.Nil && result.OwnerID == input.OwnerID {
		return result, nil
	}
	if input.ProviderID != uuid.Nil && result.ProviderID == input.ProviderID {
		return result, nil
	}

	err = fmt.Errorf("não autorizado a visualizar esta solicitação")
	return nil, err
}

// ListRequestsForOwnerUseCase lista solicitações do dono.
type ListRequestsForOwnerUseCase struct {
	requestRepo domainRequest.Repository
	logger      *slog.Logger
}

func NewListRequestsForOwnerUseCase(requestRepo domainRequest.Repository, logger *slog.Logger) *ListRequestsForOwnerUseCase {
	return &ListRequestsForOwnerUseCase{requestRepo: requestRepo, logger: logging.EnsureLogger(logger)}
}

// ListRequestsForOwnerInput filtro de listagem por dono.
type ListRequestsForOwnerInput struct {
	OwnerID uuid.UUID
	Page    int
	Limit   int
}

func (uc *ListRequestsForOwnerUseCase) Execute(ctx context.Context, input ListRequestsForOwnerInput) ([]*domainRequest.ServiceRequest, int64, error) {
	var (
		requests []*domainRequest.ServiceRequest
		total    int64
		err      error
	)
	defer logging.UseCase(ctx, uc.logger, "ListRequestsForOwnerUseCase", slog.String("owner_id", input.OwnerID.String()))(&err)

	if input.OwnerID == uuid.Nil {
		err = fmt.Errorf("ownerID é obrigatório")
		return nil, 0, err
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err = uc.requestRepo.FindByOwnerID(ctx, input.OwnerID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListRequestsForProviderUseCase lista solicitações do prestador.
type ListRequestsForProviderUseCase struct {
	requestRepo domainRequest.Repository
	logger      *slog.Logger
}

func NewListRequestsForProviderUseCase(requestRepo domainRequest.Repository, logger *slog.Logger) *ListRequestsForProviderUseCase {
	return &ListRequestsForProviderUseCase{requestRepo: requestRepo, logger: logging.EnsureLogger(logger)}
}

// ListRequestsForProviderInput filtro de listagem por prestador.
type ListRequestsForProviderInput struct {
	ProviderID uuid.UUID
	Page       int
	Limit      int
}

func (uc *ListRequestsForProviderUseCase) Execute(ctx context.Context, input ListRequestsForProviderInput) ([]*domainRequest.ServiceRequest, int64, error) {
	var (
		requests []*domainRequest.ServiceRequest
		total    int64
		err      error
	)
	defer logging.UseCase(ctx, uc.logger, "ListRequestsForProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	if input.ProviderID == uuid.Nil {
		err = fmt.Errorf("providerID é obrigatório")
		return nil, 0, err
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err = uc.requestRepo.FindByProviderID(ctx, input.ProviderID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListRequestsByStatusUseCase lista solicitações por status (uso administrativo ou dashboards).
type ListRequestsByStatusUseCase struct {
	requestRepo domainRequest.Repository
	logger      *slog.Logger
}

func NewListRequestsByStatusUseCase(requestRepo domainRequest.Repository, logger *slog.Logger) *ListRequestsByStatusUseCase {
	return &ListRequestsByStatusUseCase{requestRepo: requestRepo, logger: logging.EnsureLogger(logger)}
}

// ListRequestsByStatusInput filtro de listagem por status.
type ListRequestsByStatusInput struct {
	Status domainRequest.Status
	Page   int
	Limit  int
}

func (uc *ListRequestsByStatusUseCase) Execute(ctx context.Context, input ListRequestsByStatusInput) ([]*domainRequest.ServiceRequest, int64, error) {
	var (
		requests []*domainRequest.ServiceRequest
		total    int64
		err      error
	)
	defer logging.UseCase(ctx, uc.logger, "ListRequestsByStatusUseCase", slog.String("status", string(input.Status)))(&err)

	if input.Status == "" {
		err = fmt.Errorf("status é obrigatório")
		return nil, 0, err
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err = uc.requestRepo.FindByStatus(ctx, input.Status, page, limit)
	if err != nil {
		return nil, 0, err
	}

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
