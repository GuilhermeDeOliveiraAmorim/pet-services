package user

import (
	"context"
	"log/slog"
	"time"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ConfirmEmailVerificationUseCase confirma a verificação de email.
type ConfirmEmailVerificationUseCase struct {
	userRepo domainUser.Repository
	logger   *slog.Logger
}

func NewConfirmEmailVerificationUseCase(userRepo domainUser.Repository, logger *slog.Logger) *ConfirmEmailVerificationUseCase {
	return &ConfirmEmailVerificationUseCase{
		userRepo: userRepo,
		logger:   logging.EnsureLogger(logger),
	}
}

// ConfirmEmailVerificationInput dados para confirmar verificação.
type ConfirmEmailVerificationInput struct {
	Token string
}

// Execute valida o token e marca o email como verificado.
func (uc *ConfirmEmailVerificationUseCase) Execute(ctx context.Context, input ConfirmEmailVerificationInput) error {
	var (
		err    error
		userID string
	)
	defer logging.UseCase(ctx, uc.logger, "ConfirmEmailVerificationUseCase", slog.String("user_id", userID))(&err)

	if input.Token == "" {
		err = domainUser.ErrEmailVerificationTokenInvalid
		return err
	}

	// Busca o token de verificação
	verificationToken, err := uc.userRepo.FindEmailVerificationToken(ctx, input.Token)
	if err != nil {
		err = domainUser.ErrEmailVerificationTokenInvalid
		return err
	}

	// Verifica se o token expirou
	if time.Now().After(verificationToken.ExpiresAt) {
		err = domainUser.ErrEmailVerificationTokenInvalid
		return err
	}

	// Busca o usuário
	user, err := uc.userRepo.FindByID(ctx, verificationToken.UserID)
	if err != nil {
		return err
	}
	userID = user.ID.String()

	// Verifica se já está verificado
	if user.EmailVerified {
		err = domainUser.ErrEmailAlreadyVerified
		return err
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
