package usecases

import (
	"context"
	"errors"
	"time"

	"pet-services-api/internal/auth"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type LoginUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type LoginUserOutput struct {
	User         *entities.User              `json:"user"`
	AccessToken  string                      `json:"access_token"`
	RefreshToken string                      `json:"refresh_token"`
	ExpiresIn    int64                       `json:"expires_in"`
	Problems     []exceptions.ProblemDetails `json:"problems,omitempty"`
}

type LoginUserUseCase struct {
	userRepository         entities.UserRepository
	refreshTokenRepository entities.RefreshTokenRepository
	storage                storage.ObjectStorage
	logger                 logging.LoggerInterface
}

func NewLoginUserUseCase(userRepo entities.UserRepository, refreshRepo entities.RefreshTokenRepository, storageService storage.ObjectStorage, logger logging.LoggerInterface) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository:         userRepo,
		refreshTokenRepository: refreshRepo,
		storage:                storageService,
		logger:                 logger,
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
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o email informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	login := entities.Login{Email: user.Login.Email, Password: user.Login.Password}
	if !login.CompareAndDecrypt(user.Login.Password, input.Password) {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Credenciais inválidas", errors.New("Email ou senha incorretos"))
	}

	if !user.Active {
		return nil, uc.logger.LogForbidden(ctx, from, "Conta desativada", errors.New("Sua conta foi desativada. Entre em contato com o suporte para reativar"))
	}

	if !user.EmailVerified {
		return nil, uc.logger.LogForbidden(ctx, from, "Email não verificado", errors.New("Verifique seu email antes de fazer login. Utilize a opção de reenviar email de verificação"))
	}

	if err := signUserPhotos(ctx, uc.storage, user); err != nil {
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
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
