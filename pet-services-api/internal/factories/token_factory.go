package factories

import (
	"pet-services-api/internal/mail"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"
	"time"

	"gorm.io/gorm"
)

type TokenFactory struct {
	LoginUser               *usecases.LoginUserUseCase
	Logout                  *usecases.LogoutUseCase
	RequestPasswordReset    *usecases.RequestPasswordResetUseCase
	ResetPassword           *usecases.ResetPasswordUseCase
	ResendVerificationEmail *usecases.ResendVerificationEmailUseCase
	VerifyEmail             *usecases.VerifyEmailUseCase
}

func NewTokenFactory(db *gorm.DB, mailService mail.EmailService, ttl time.Duration) *TokenFactory {
	userRepo := repository_impl.NewUserRepository(db)
	tokenRepo := repository_impl.NewRefreshTokenRepository(db)

	return &TokenFactory{
		LoginUser:               usecases.NewLoginUserUseCase(userRepo, tokenRepo),
		Logout:                  usecases.NewLogoutUseCase(tokenRepo),
		RequestPasswordReset:    usecases.NewRequestPasswordResetUseCase(userRepo, tokenRepo, mailService, ttl),
		ResetPassword:           usecases.NewResetPasswordUseCase(userRepo, tokenRepo),
		ResendVerificationEmail: usecases.NewResendVerificationEmailUseCase(userRepo, tokenRepo, mailService, ttl),
		VerifyEmail:             usecases.NewVerifyEmailUseCase(userRepo, tokenRepo),
	}
}
