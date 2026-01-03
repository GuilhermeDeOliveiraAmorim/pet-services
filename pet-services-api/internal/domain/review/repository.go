package review

import (
	"context"

	"github.com/google/uuid"
)

// Repository define a interface para persistência de avaliações.
// Implementações ficam na camada de infraestrutura (ex.: Postgres, cache, etc.).
type Repository interface {
	// Create cria uma nova avaliação.
	Create(ctx context.Context, review *Review) error

	// FindByID busca uma avaliação por ID.
	FindByID(ctx context.Context, id uuid.UUID) (*Review, error)

	// FindByRequestID busca avaliação vinculada a uma solicitação.
	FindByRequestID(ctx context.Context, requestID uuid.UUID) (*Review, error)

	// ExistsByRequestID verifica se já existe avaliação para a solicitação.
	ExistsByRequestID(ctx context.Context, requestID uuid.UUID) (bool, error)

	// ListByProvider lista avaliações de um prestador com paginação.
	ListByProvider(ctx context.Context, providerID uuid.UUID, page, limit int) ([]*Review, int64, error)

	// Delete remove uma avaliação.
	Delete(ctx context.Context, id uuid.UUID) error
}
