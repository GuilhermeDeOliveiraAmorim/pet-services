package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// SignupUseCase registra um novo usuário e retorna tokens.
type SignupUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	tokenService   domainAuth.TokenService
	passwordHasher domainAuth.PasswordHasher
	logger         logging.LoggerService
}

func NewSignupUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, passwordHasher domainAuth.PasswordHasher, logger logging.LoggerService) *SignupUseCase {
	return &SignupUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// SignupInput dados para criação de usuário.
type SignupInput struct {
	Email    string
	Name     string
	Phone    string
	Password string
	Type     domainUser.UserType
}

// SignupOutput resultado com usuário e tokens.
type SignupOutput struct {
	UserID           uuid.UUID
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  int64
	RefreshExpiresAt int64
	UserType         domainUser.UserType
}

const SIGNUP_USECASE = "SIGNUP_USECASE"

// Execute cria usuário, hash de senha, tokens e registra refresh, seguindo padrão de erros e logging.
func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (*SignupOutput, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    SIGNUP_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	email := strings.TrimSpace(strings.ToLower(input.Email))
	name := strings.TrimSpace(input.Name)
	password := strings.TrimSpace(input.Password)

	if password == "" || len(password) < 6 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SIGNUP_USECASE,
			Message: "Senha inválida",
			Error:   errors.New("senha inválida"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Senha inválida",
			Status: exceptions.RFC400_CODE,
			Detail: "A senha deve ter pelo menos 6 caracteres.",
		}}
	}

	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SIGNUP_USECASE,
			Message: "Erro ao verificar e-mail existente",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao verificar e-mail existente",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}
	if exists {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    SIGNUP_USECASE,
			Message: "E-mail já cadastrado",
			Error:   errors.New("e-mail já cadastrado"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "E-mail já cadastrado",
			Status: exceptions.RFC409_CODE,
			Detail: "O e-mail informado já está cadastrado.",
		}}
	}

	user, err := domainUser.NewUser(email, name, input.Phone, input.Type)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SIGNUP_USECASE,
			Message: "Dados de usuário inválidos",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Dados de usuário inválidos",
			Status: exceptions.RFC400_CODE,
			Detail: err.Error(),
		}}
	}

	hashed, err := uc.passwordHasher.Hash(password)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SIGNUP_USECASE,
			Message: "Erro ao gerar hash da senha",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao gerar hash da senha",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}
	user.SetPassword(hashed)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SIGNUP_USECASE,
			Message: "Erro ao criar usuário",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao criar usuário",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	tokens, err := uc.tokenService.GenerateTokens(user.ID, user.Type)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SIGNUP_USECASE,
			Message: "Erro ao gerar tokens",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao gerar tokens",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	rt := domainAuth.NewRefreshToken(user.ID, tokens.RefreshExpiresAt)
	// usa o ID gerado pelo token service para rastrear o mesmo token
	if tokens.RefreshID != uuid.Nil {
		rt.ID = tokens.RefreshID
	}
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SIGNUP_USECASE,
			Message: "Falha ao salvar refresh token",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar refresh token",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    SIGNUP_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	result := &SignupOutput{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserType:         user.Type,
	}

	return result, nil
}
