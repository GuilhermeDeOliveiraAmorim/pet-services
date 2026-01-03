package user

import (
	"context"
	"time"

	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ConfirmEmailVerificationUseCase confirma a verificação de email.
type ConfirmEmailVerificationUseCase struct {
	userRepo domainUser.Repository
}

func NewConfirmEmailVerificationUseCase(userRepo domainUser.Repository) *ConfirmEmailVerificationUseCase {
	return &ConfirmEmailVerificationUseCase{
		userRepo: userRepo,
	}
}

// ConfirmEmailVerificationInput dados para confirmar verificação.
type ConfirmEmailVerificationInput struct {
	Token string
}

// Execute valida o token e marca o email como verificado.
func (uc *ConfirmEmailVerificationUseCase) Execute(ctx context.Context, input ConfirmEmailVerificationInput) error {
	if input.Token == "" {
		return domainUser.ErrEmailVerificationTokenInvalid
	}

	// Busca o token de verificação
	verificationToken, err := uc.userRepo.FindEmailVerificationToken(ctx, input.Token)
	if err != nil {
		return domainUser.ErrEmailVerificationTokenInvalid
	}

	// Verifica se o token expirou
	if time.Now().After(verificationToken.ExpiresAt) {
		return domainUser.ErrEmailVerificationTokenInvalid
	}

	// Busca o usuário
	user, err := uc.userRepo.FindByID(ctx, verificationToken.UserID)
	if err != nil {
		return err
	}

	// Verifica se já está verificado
	if user.EmailVerified {
		return domainUser.ErrEmailAlreadyVerified
	}

	// Marca o email como verificado
	user.VerifyEmail()

	// Atualiza o usuário
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Remove o token de verificação usado
	if err := uc.userRepo.MarkEmailVerificationTokenAsUsed(ctx, verificationToken.ID); err != nil {
		// Log error mas não falha a operação, pois a verificação já foi feita
		return nil
	}

	return nil
}
