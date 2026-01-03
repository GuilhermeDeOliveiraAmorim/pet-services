package request

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// GetRequestStatusUseCase retorna uma solicitação garantindo vínculo com dono ou prestador.
type GetRequestStatusUseCase struct {
	requestRepo domainRequest.Repository
}

func NewGetRequestStatusUseCase(requestRepo domainRequest.Repository) *GetRequestStatusUseCase {
	return &GetRequestStatusUseCase{requestRepo: requestRepo}
}

// GetRequestStatusInput entrada para buscar status.
type GetRequestStatusInput struct {
	RequestID  uuid.UUID
	OwnerID    uuid.UUID
	ProviderID uuid.UUID
}

func (uc *GetRequestStatusUseCase) Execute(ctx context.Context, input GetRequestStatusInput) (*domainRequest.ServiceRequest, error) {
	if input.RequestID == uuid.Nil {
		return nil, fmt.Errorf("requestID é obrigatório")
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	// Autorização simples: precisa ser dono ou prestador vinculado.
	if input.OwnerID != uuid.Nil && req.OwnerID == input.OwnerID {
		return req, nil
	}
	if input.ProviderID != uuid.Nil && req.ProviderID == input.ProviderID {
		return req, nil
	}

	return nil, fmt.Errorf("não autorizado a visualizar esta solicitação")
}

// ListRequestsForOwnerUseCase lista solicitações do dono.
type ListRequestsForOwnerUseCase struct {
	requestRepo domainRequest.Repository
}

func NewListRequestsForOwnerUseCase(requestRepo domainRequest.Repository) *ListRequestsForOwnerUseCase {
	return &ListRequestsForOwnerUseCase{requestRepo: requestRepo}
}

// ListRequestsForOwnerInput filtro de listagem por dono.
type ListRequestsForOwnerInput struct {
	OwnerID uuid.UUID
	Page    int
	Limit   int
}

func (uc *ListRequestsForOwnerUseCase) Execute(ctx context.Context, input ListRequestsForOwnerInput) ([]*domainRequest.ServiceRequest, int64, error) {
	if input.OwnerID == uuid.Nil {
		return nil, 0, fmt.Errorf("ownerID é obrigatório")
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByOwnerID(ctx, input.OwnerID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListRequestsForProviderUseCase lista solicitações do prestador.
type ListRequestsForProviderUseCase struct {
	requestRepo domainRequest.Repository
}

func NewListRequestsForProviderUseCase(requestRepo domainRequest.Repository) *ListRequestsForProviderUseCase {
	return &ListRequestsForProviderUseCase{requestRepo: requestRepo}
}

// ListRequestsForProviderInput filtro de listagem por prestador.
type ListRequestsForProviderInput struct {
	ProviderID uuid.UUID
	Page       int
	Limit      int
}

func (uc *ListRequestsForProviderUseCase) Execute(ctx context.Context, input ListRequestsForProviderInput) ([]*domainRequest.ServiceRequest, int64, error) {
	if input.ProviderID == uuid.Nil {
		return nil, 0, fmt.Errorf("providerID é obrigatório")
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByProviderID(ctx, input.ProviderID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListRequestsByStatusUseCase lista solicitações por status (uso administrativo ou dashboards).
type ListRequestsByStatusUseCase struct {
	requestRepo domainRequest.Repository
}

func NewListRequestsByStatusUseCase(requestRepo domainRequest.Repository) *ListRequestsByStatusUseCase {
	return &ListRequestsByStatusUseCase{requestRepo: requestRepo}
}

// ListRequestsByStatusInput filtro de listagem por status.
type ListRequestsByStatusInput struct {
	Status domainRequest.Status
	Page   int
	Limit  int
}

func (uc *ListRequestsByStatusUseCase) Execute(ctx context.Context, input ListRequestsByStatusInput) ([]*domainRequest.ServiceRequest, int64, error) {
	if input.Status == "" {
		return nil, 0, fmt.Errorf("status é obrigatório")
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	requests, total, err := uc.requestRepo.FindByStatus(ctx, input.Status, page, limit)
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
