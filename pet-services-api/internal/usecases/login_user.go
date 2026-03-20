package usecases

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"pet-services-api/internal/auth"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/storage"
)

var blockedLoginAlertLimiter = newLoginBlockedAlertLimiter(30 * time.Minute)

type loginBlockedAlertLimiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
}

func newLoginBlockedAlertLimiter(cooldown time.Duration) *loginBlockedAlertLimiter {
	return &loginBlockedAlertLimiter{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

func (l *loginBlockedAlertLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	last, ok := l.lastSent[key]
	if ok && now.Sub(last) < l.cooldown {
		return false
	}

	l.lastSent[key] = now
	return true
}

type LoginUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type LoginUserOutput struct {
	User         *UserOutput                 `json:"user"`
	AccessToken  string                      `json:"access_token"`
	RefreshToken string                      `json:"refresh_token"`
	ExpiresIn    int64                       `json:"expires_in"`
	Problems     []exceptions.ProblemDetails `json:"problems,omitempty"`
}

type LoginUserUseCase struct {
	userRepository         entities.UserRepository
	refreshTokenRepository entities.RefreshTokenRepository
	storage                storage.ObjectStorage
	emailService           mail.EmailService
	logger                 logging.LoggerInterface
}

func NewLoginUserUseCase(userRepo entities.UserRepository, refreshRepo entities.RefreshTokenRepository, storageService storage.ObjectStorage, emailService mail.EmailService, logger logging.LoggerInterface) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository:         userRepo,
		refreshTokenRepository: refreshRepo,
		storage:                storageService,
		emailService:           emailService,
		logger:                 logger,
	}
}

func (uc *LoginUserUseCase) notifyBlockedLoginAttempt(ctx context.Context, user *entities.User, reason string) {
	const from = "LoginUserUseCase.notifyBlockedLoginAttempt"

	key := strings.ToLower(strings.TrimSpace(user.Login.Email)) + ":" + reason
	if !blockedLoginAlertLimiter.allow(key, time.Now()) {
		return
	}

	if err := uc.emailService.SendLoginBlockedAlertEmail(user.Login.Email, user.Name, reason); err != nil {
		uc.logger.LogWarning(ctx, from, "Falha ao enviar alerta de login bloqueado", err)
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, []exceptions.ProblemDetails) {
	const from = "LoginUserUseCase.Execute"

	if input.Email == "" || input.Password == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Credenciais ausentes", errors.New("Email e senha são obrigatórios"))
	}

	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogUnauthorized(ctx, from, "Credenciais inválidas", errors.New("Credenciais inválidas"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	login := entities.Login{Email: user.Login.Email, Password: user.Login.Password}
	if !login.CompareAndDecrypt(user.Login.Password, input.Password) {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Credenciais inválidas", errors.New("Credenciais inválidas"))
	}

	if !user.Active {
		uc.notifyBlockedLoginAttempt(ctx, user, "Conta desativada")
		return nil, uc.logger.LogForbidden(ctx, from, "Conta desativada", errors.New("Sua conta foi desativada. Entre em contato com o suporte para reativar"))
	}

	if !user.EmailVerified {
		uc.notifyBlockedLoginAttempt(ctx, user, "Email nao verificado")
		return nil, uc.logger.LogForbidden(ctx, from, "Email não verificado", errors.New("Verifique seu email antes de fazer login. Utilize a opção de reenviar email de verificação"))
	}

	if err := storage.SignUserPhotos(ctx, uc.storage, user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URLs das fotos", err)
	}

	jwtSvc, err := auth.NewJWTServiceFromEnv()
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao carregar configuração JWT", err)
	}

	accessToken, err := jwtSvc.GenerateAccessToken(user.ID, user.Login.Email, user.UserType)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar token de acesso", err)
	}

	refreshToken, err := jwtSvc.GenerateRefreshToken(user.ID, user.Login.Email, user.UserType)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar token de atualização", err)
	}

	claims, err := jwtSvc.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao validar token de acesso", err)
	}

	expiresIn := max(int64(time.Until(claims.ExpiresAt.Time).Seconds()), 0)

	refreshClaims, err := jwtSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao validar token de atualização", err)
	}

	refreshEntity, problems := entities.NewRefreshToken(user.ID, refreshToken, refreshClaims.ExpiresAt.Time, input.UserAgent, input.IP)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Token de atualização inválido", problems)
		return nil, problems
	}

	if err := uc.refreshTokenRepository.RevokeAllByUserID(user.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar tokens de atualização", err)
	}

	if err := uc.refreshTokenRepository.Create(refreshEntity); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar token de atualização", err)
	}

	return &LoginUserOutput{
		User:         NewUserOutput(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
