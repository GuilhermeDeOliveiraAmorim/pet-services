package request

import (
	"context"

	"github.com/google/uuid"
)

// Repository define a interface para persistência de solicitações
type Repository interface {
	// Create cria uma nova solicitação
	Create(ctx context.Context, request *ServiceRequest) error

	// FindByID busca uma solicitação por ID
	FindByID(ctx context.Context, id uuid.UUID) (*ServiceRequest, error)

	// Update atualiza uma solicitação
	Update(ctx context.Context, request *ServiceRequest) error

	// Delete remove uma solicitação
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByOwnerID busca solicitações por ID do dono
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID, page, limit int) ([]*ServiceRequest, int64, error)

	// FindByProviderID busca solicitações por ID do prestador
	FindByProviderID(ctx context.Context, providerID uuid.UUID, page, limit int) ([]*ServiceRequest, int64, error)

	// FindByStatus busca solicitações por status
	FindByStatus(ctx context.Context, status Status, page, limit int) ([]*ServiceRequest, int64, error)
}
