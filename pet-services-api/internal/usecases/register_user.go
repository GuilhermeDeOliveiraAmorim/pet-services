package usecases

import (
	"context"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"time"

	"github.com/oklog/ulid/v2"
)

type RegisterUserInput struct {
	Name     string         `json:"name"`
	UserType string         `json:"user_type"`
	Login    entities.Login `json:"login"`
	Phone    entities.Phone `json:"phone"`
}

type RegisterUserOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type RegisterUserUseCase struct {
	userRepository        entities.UserRepository
	verifyTokenRepository entities.RefreshTokenRepository
	emailService          mail.EmailService
	ttl                   time.Duration
	logger                logging.LoggerInterface
}

func NewRegisterUserUseCase(userRepository entities.UserRepository, verifyTokenRepository entities.RefreshTokenRepository, emailService mail.EmailService, ttl time.Duration, logger logging.LoggerInterface) *RegisterUserUseCase {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &RegisterUserUseCase{
		userRepository:        userRepository,
		verifyTokenRepository: verifyTokenRepository,
		emailService:          emailService,
		ttl:                   ttl,
		logger:                logger,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, []exceptions.ProblemDetails) {
	const from = "RegisterUserUseCase.Execute"

	exists, err := uc.userRepository.ExistsByEmail(input.Login.Email)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar email", err)
	}

	if exists {
		return &RegisterUserOutput{
			Message: "Usuário registrado com sucesso",
			Detail:  "O usuário foi criado com sucesso. Verifique seu email para ativar a conta",
		}, nil
	}

	var problems []exceptions.ProblemDetails

	var login *entities.Login
	loginResult, errs := entities.NewLogin(input.Login.Email, input.Login.Password)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Login inválido", errs)
		problems = append(problems, errs...)
	}

	if len(errs) == 0 {
		login = loginResult
	}

	var phone *entities.Phone
	phoneResult, errs := entities.NewPhone(input.Phone.CountryCode, input.Phone.AreaCode, input.Phone.Number)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Telefone inválido", errs)
		problems = append(problems, errs...)
	}

	if len(errs) == 0 {
		phone = phoneResult
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	userType := input.UserType
	if userType == "" {
		userType = entities.UserTypes.Owner
	}

	user, errs := entities.NewIncompleteUser(
		input.Name,
		userType,
		*login,
		*phone,
	)

	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Usuário inválido", errs)
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	if err := user.Login.EncryptPassword(); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	if err := uc.userRepository.Create(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar usuário", err)
	}

	tokenStr := ulid.Make().String()
	expiresAt := time.Now().Add(uc.ttl)

	if err := uc.verifyTokenRepository.RevokeAllPasswordResetByUserID(user.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar tokens anteriores", err)
	}

	verifyToken := &entities.PasswordResetToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := uc.verifyTokenRepository.CreatePasswordReset(verifyToken); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar token de verificação", err)
	}

	if err := uc.emailService.SendVerificationEmail(user.Login.Email, tokenStr); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email de verificação", err)
	}

	return &RegisterUserOutput{
		Message: "Usuário registrado com sucesso",
		Detail:  "O usuário foi criado com sucesso. Verifique seu email para ativar a conta",
	}, nil
}
