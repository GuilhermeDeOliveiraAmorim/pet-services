package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository define a interface para persistência de usuários.
// Implementações ficam na camada de infraestrutura (ex.: Postgres, cache, etc.).
type Repository interface {
	// Create cria um novo usuário.
	Create(ctx context.Context, user *User) error

	// FindByID busca um usuário por ID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail busca um usuário pelo email.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Update atualiza os dados de um usuário.
	Update(ctx context.Context, user *User) error

	// Delete remove um usuário.
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByEmail verifica se já existe usuário com o email.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
