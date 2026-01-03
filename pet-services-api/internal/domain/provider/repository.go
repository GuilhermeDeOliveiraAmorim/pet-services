package provider

import (
	"context"

	"github.com/google/uuid"
)

// Repository define a interface para persistência de prestadores.
// Implementações concretas ficam na camada de infraestrutura (ex.: Postgres, Elasticsearch, cache, etc.).
type Repository interface {
	// Create cria um novo prestador.
	Create(ctx context.Context, provider *Provider) error

	// FindByID busca um prestador pelo ID.
	FindByID(ctx context.Context, id uuid.UUID) (*Provider, error)

	// FindByUserID busca o prestador vinculado a um usuário.
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Provider, error)

	// Update atualiza os dados do prestador.
	Update(ctx context.Context, provider *Provider) error

	// Delete remove um prestador.
	Delete(ctx context.Context, id uuid.UUID) error

	// List lista prestadores paginados.
	List(ctx context.Context, page, limit int) ([]*Provider, int64, error)

	// FindActiveByLocation busca prestadores ativos próximos a uma localização (raio em km).
	FindActiveByLocation(ctx context.Context, latitude, longitude, radiusKM float64, page, limit int) ([]*Provider, int64, error)

	// ExistsByUserID verifica se já existe prestador para um usuário.
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}
