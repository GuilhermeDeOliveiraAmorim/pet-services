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

	// CreatePasswordResetToken cria um token de reset de senha.
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error

	// FindPasswordResetToken busca token de reset por valor.
	FindPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)

	// MarkPasswordResetTokenAsUsed marca o token como utilizado.
	MarkPasswordResetTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error

	// CreateEmailVerificationToken cria um token de verificação de email.
	CreateEmailVerificationToken(ctx context.Context, token *EmailVerificationToken) error

	// FindEmailVerificationToken busca token de verificação por valor.
	FindEmailVerificationToken(ctx context.Context, token string) (*EmailVerificationToken, error)

	// MarkEmailVerificationTokenAsUsed marca o token como utilizado.
	MarkEmailVerificationTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error
}
