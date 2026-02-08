package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type RegisterUserInput struct {
	Name     string           `json:"name"`
	UserType string           `json:"user_type"`
	Login    entities.Login   `json:"login"`
	Phone    entities.Phone   `json:"phone"`
	Address  entities.Address `json:"address"`
}

type RegisterUserOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type RegisterUserUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewRegisterUserUseCase(userRepository entities.UserRepository, logger logging.LoggerInterface) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, []exceptions.ProblemDetails) {
	const from = "RegisterUserUseCase.Execute"

	if input.UserType == "admin" {
		return nil, uc.logger.LogForbidden(ctx, from, "Tipo de usuário não permitido", errors.New("Não é permitido criar usuários do tipo admin por este endpoint"))
	}

	exists, err := uc.userRepository.ExistsByEmail(input.Login.Email)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar email", err)
	}

	if exists {
		return nil, uc.logger.LogConflict(ctx, from, "Email já cadastrado", errors.New("O email informado já está em uso por outro usuário"))
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

	var location *entities.Location
	locationResult, errs := entities.NewLocation(input.Address.Location.Latitude, input.Address.Location.Longitude)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Localização inválida", errs)
		problems = append(problems, errs...)
	}
	if len(errs) == 0 {
		location = locationResult
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	address, errs := entities.NewAddress(
		input.Address.Street,
		input.Address.Number,
		input.Address.Neighborhood,
		input.Address.City,
		input.Address.ZipCode,
		input.Address.State,
		input.Address.Country,
		input.Address.Complement,
		*location,
	)
	if len(errs) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Endereço inválido", errs)
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Problemas de validação", problems)
		return nil, problems
	}

	user, errs := entities.NewUser(
		input.Name,
		input.UserType,
		*login,
		*phone,
		*address,
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

	return &RegisterUserOutput{
		Message: "Usuário registrado com sucesso",
		Detail:  "O usuário foi criado e registrado no sistema com sucesso",
	}, nil
}
