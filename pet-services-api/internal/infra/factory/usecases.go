package factory

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/guilherme/pet-services-api/internal/application/auth"
	"github.com/guilherme/pet-services-api/internal/application/provider"
	"github.com/guilherme/pet-services-api/internal/application/request"
	"github.com/guilherme/pet-services-api/internal/application/review"
	"github.com/guilherme/pet-services-api/internal/application/user"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	providerdom "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
	gormrepo "github.com/guilherme/pet-services-api/internal/infra/repository/gorm"
)

// Config encapsula dependências externas necessárias para criar os casos de uso.
type Config struct {
	DB                   *gorm.DB
	TokenService         domainAuth.TokenService
	PasswordHasher       domainAuth.PasswordHasher
	EmailService         domainUser.EmailService
	PasswordResetBaseURL string
	EmailVerifyBaseURL   string
	Logger               *slog.Logger
}

// UseCases agrupa fábricas de casos de uso por contexto.
type UseCases struct {
	Auth         AuthUseCases
	User         UserUseCases
	Provider     ProviderUseCases
	Request      RequestUseCases
	Review       ReviewUseCases
	ProviderRepo providerdom.Repository
}

// AuthUseCases contém os casos de uso do contexto de autenticação.
type AuthUseCases struct {
	Login   *auth.LoginUseCase
	Signup  *auth.SignupUseCase
	Refresh *auth.RefreshTokenUseCase
	Logout  *auth.LogoutUseCase
}

// UserUseCases contém os casos de uso relacionados ao usuário.
type UserUseCases struct {
	GetProfile               *user.GetProfileUseCase
	UpdateProfile            *user.UpdateProfileUseCase
	ChangePassword           *user.ChangePasswordUseCase
	RequestPasswordReset     *user.RequestPasswordResetUseCase
	ConfirmPasswordReset     *user.ConfirmPasswordResetUseCase
	RequestEmailVerification *user.RequestEmailVerificationUseCase
	ConfirmEmailVerification *user.ConfirmEmailVerificationUseCase
	DeleteAccount            *user.DeleteAccountUseCase
}

// ProviderUseCases contém os casos de uso relacionados ao prestador.
type ProviderUseCases struct {
	Register           *provider.RegisterProviderUseCase
	UpdateProfile      *provider.UpdateProviderProfileUseCase
	Activate           *provider.ActivateProviderUseCase
	Deactivate         *provider.DeactivateProviderUseCase
	Approve            *provider.ApproveProviderUseCase
	Reject             *provider.RejectProviderUseCase
	AddService         *provider.AddServiceUseCase
	UpdateService      *provider.UpdateServiceUseCase
	RemoveService      *provider.RemoveServiceUseCase
	AddPhoto           *provider.AddPhotoUseCase
	RemovePhoto        *provider.RemovePhotoUseCase
	ReorderPhotos      *provider.ReorderPhotosUseCase
	UpdateWorkingHours *provider.UpdateWorkingHoursUseCase
	UpdateDaySchedule  *provider.UpdateDayScheduleUseCase
	ListByLocation     *provider.ListProvidersByLocationUseCase
}

// RequestUseCases contém os casos de uso do fluxo de solicitações de serviço.
type RequestUseCases struct {
	Create          *request.CreateRequestUseCase
	Accept          *request.AcceptRequestUseCase
	Reject          *request.RejectRequestUseCase
	Complete        *request.CompleteRequestUseCase
	Cancel          *request.CancelRequestUseCase
	GetStatus       *request.GetRequestStatusUseCase
	ListForOwner    *request.ListRequestsForOwnerUseCase
	ListForProvider *request.ListRequestsForProviderUseCase
	ListByStatus    *request.ListRequestsByStatusUseCase
}

// ReviewUseCases contém os casos de uso de avaliações.
type ReviewUseCases struct {
	Submit          *review.SubmitReviewUseCase
	ListForProvider *review.ListReviewsForProviderUseCase
}

// NewUseCases cria todas as dependências e retorna o agregador de casos de uso.
func NewUseCases(cfg Config) UseCases {
	userRepo := gormrepo.NewUserRepository(cfg.DB)
	providerRepo := gormrepo.NewProviderRepository(cfg.DB)
	requestRepo := gormrepo.NewRequestRepository(cfg.DB)
	reviewRepo := gormrepo.NewReviewRepository(cfg.DB)
	refreshRepo := gormrepo.NewRefreshTokenRepository(cfg.DB)

	return UseCases{
		ProviderRepo: providerRepo,
		Auth: AuthUseCases{
			Login:   auth.NewLoginUseCase(userRepo, refreshRepo, cfg.TokenService, cfg.PasswordHasher, cfg.Logger),
			Signup:  auth.NewSignupUseCase(userRepo, refreshRepo, cfg.TokenService, cfg.PasswordHasher, cfg.Logger),
			Refresh: auth.NewRefreshTokenUseCase(refreshRepo, cfg.TokenService, cfg.Logger),
			Logout:  auth.NewLogoutUseCase(refreshRepo, cfg.TokenService, cfg.Logger),
		},
		User: UserUseCases{
			GetProfile:               user.NewGetProfileUseCase(userRepo, cfg.Logger),
			UpdateProfile:            user.NewUpdateProfileUseCase(userRepo, cfg.Logger),
			ChangePassword:           user.NewChangePasswordUseCase(userRepo, refreshRepo, cfg.PasswordHasher, cfg.Logger),
			RequestPasswordReset:     user.NewRequestPasswordResetUseCase(userRepo, cfg.EmailService, cfg.PasswordResetBaseURL, cfg.Logger),
			ConfirmPasswordReset:     user.NewConfirmPasswordResetUseCase(userRepo, refreshRepo, cfg.PasswordHasher, cfg.Logger),
			RequestEmailVerification: user.NewRequestEmailVerificationUseCase(userRepo, cfg.EmailService, cfg.EmailVerifyBaseURL, cfg.Logger),
			ConfirmEmailVerification: user.NewConfirmEmailVerificationUseCase(userRepo, cfg.Logger),
			DeleteAccount:            user.NewDeleteAccountUseCase(userRepo, refreshRepo, cfg.Logger),
		},
		Provider: ProviderUseCases{
			Register:           provider.NewRegisterProviderUseCase(providerRepo, userRepo, cfg.Logger),
			UpdateProfile:      provider.NewUpdateProviderProfileUseCase(providerRepo, cfg.Logger),
			Activate:           provider.NewActivateProviderUseCase(providerRepo, cfg.Logger),
			Deactivate:         provider.NewDeactivateProviderUseCase(providerRepo, cfg.Logger),
			Approve:            provider.NewApproveProviderUseCase(providerRepo, cfg.Logger),
			Reject:             provider.NewRejectProviderUseCase(providerRepo, cfg.Logger),
			AddService:         provider.NewAddServiceUseCase(providerRepo, cfg.Logger),
			UpdateService:      provider.NewUpdateServiceUseCase(providerRepo, cfg.Logger),
			RemoveService:      provider.NewRemoveServiceUseCase(providerRepo, cfg.Logger),
			AddPhoto:           provider.NewAddPhotoUseCase(providerRepo, cfg.Logger),
			RemovePhoto:        provider.NewRemovePhotoUseCase(providerRepo, cfg.Logger),
			ReorderPhotos:      provider.NewReorderPhotosUseCase(providerRepo, cfg.Logger),
			UpdateWorkingHours: provider.NewUpdateWorkingHoursUseCase(providerRepo, cfg.Logger),
			UpdateDaySchedule:  provider.NewUpdateDayScheduleUseCase(providerRepo, cfg.Logger),
			ListByLocation:     provider.NewListProvidersByLocationUseCase(providerRepo, cfg.Logger),
		},
		Request: RequestUseCases{
			Create:          request.NewCreateRequestUseCase(requestRepo, providerRepo, cfg.Logger),
			Accept:          request.NewAcceptRequestUseCase(requestRepo, providerRepo, cfg.Logger),
			Reject:          request.NewRejectRequestUseCase(requestRepo, providerRepo, cfg.Logger),
			Complete:        request.NewCompleteRequestUseCase(requestRepo, providerRepo, cfg.Logger),
			Cancel:          request.NewCancelRequestUseCase(requestRepo, cfg.Logger),
			GetStatus:       request.NewGetRequestStatusUseCase(requestRepo, cfg.Logger),
			ListForOwner:    request.NewListRequestsForOwnerUseCase(requestRepo, cfg.Logger),
			ListForProvider: request.NewListRequestsForProviderUseCase(requestRepo, cfg.Logger),
			ListByStatus:    request.NewListRequestsByStatusUseCase(requestRepo, cfg.Logger),
		},
		Review: ReviewUseCases{
			Submit:          review.NewSubmitReviewUseCase(reviewRepo, requestRepo, providerRepo, cfg.Logger),
			ListForProvider: review.NewListReviewsForProviderUseCase(reviewRepo, cfg.Logger),
		},
	}
}
