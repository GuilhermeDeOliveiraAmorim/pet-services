package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"
	"time"

	"gorm.io/gorm"
)

type TokenFactory struct {
	LoginUser               *usecases.LoginUserUseCase
	RefreshToken            *usecases.RefreshTokenUseCase
	Logout                  *usecases.LogoutUseCase
	RequestPasswordReset    *usecases.RequestPasswordResetUseCase
	ResetPassword           *usecases.ResetPasswordUseCase
	ResendVerificationEmail *usecases.ResendVerificationEmailUseCase
	VerifyEmail             *usecases.VerifyEmailUseCase
	Logger                  logging.LoggerInterface
}

func NewTokenFactory(db *gorm.DB, mailService mail.EmailService, ttl time.Duration, storageService storage.ObjectStorage, logger logging.LoggerInterface) *TokenFactory {
	userRepo := repository_impl.NewUserRepository(db)
	tokenRepo := repository_impl.NewRefreshTokenRepository(db)

	return &TokenFactory{
		LoginUser:               usecases.NewLoginUserUseCase(userRepo, tokenRepo, storageService, mailService, logger),
		RefreshToken:            usecases.NewRefreshTokenUseCase(userRepo, tokenRepo, logger),
		Logout:                  usecases.NewLogoutUseCase(tokenRepo, logger),
		RequestPasswordReset:    usecases.NewRequestPasswordResetUseCase(userRepo, tokenRepo, mailService, ttl, logger),
		ResetPassword:           usecases.NewResetPasswordUseCase(userRepo, tokenRepo, mailService, logger),
		ResendVerificationEmail: usecases.NewResendVerificationEmailUseCase(userRepo, tokenRepo, mailService, ttl, logger),
		VerifyEmail:             usecases.NewVerifyEmailUseCase(userRepo, tokenRepo, mailService, logger),
		Logger:                  logger,
	}
}
