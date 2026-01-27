package usecases

import (
	"context"
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
}

func NewRegisterUserUseCase(userRepository entities.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, []exceptions.ProblemDetails) {
	const from = "RegisterUserUseCase.Execute"

	if input.UserType == "admin" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
				Title:  "Tipo de usuário não permitido",
				Detail: "Não é permitido criar usuários do tipo admin por este endpoint",
			}),
		}
	}

	exists, err := uc.userRepository.ExistsByEmail(input.Login.Email)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao verificar email", err)
	}

	if exists {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.Conflict, exceptions.ErrorMessage{
				Title:  "Email já cadastrado",
				Detail: "O email informado já está em uso por outro usuário",
			}),
		}
	}

	var problems []exceptions.ProblemDetails

	login, errs := entities.NewLogin(input.Login.Email, input.Login.Password)
	if len(errs) > 0 {
		problems = append(problems, errs...)
	}

	if err := login.EncryptPassword(); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	phone, errs := entities.NewPhone(input.Phone.CountryCode, input.Phone.AreaCode, input.Phone.Number)
	if len(errs) > 0 {
		problems = append(problems, errs...)
	}

	location, errs := entities.NewLocation(input.Address.Location.Latitude, input.Address.Location.Longitude)
	if len(errs) > 0 {
		problems = append(problems, errs...)
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
		problems = append(problems, errs...)
	}

	user, errs := entities.NewUser(
		input.Name,
		input.UserType,
		*login,
		*phone,
		*address,
	)
	if len(errs) > 0 {
		problems = append(problems, errs...)
	}

	if len(problems) > 0 {
		return nil, problems
	}

	if err := uc.userRepository.Create(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criar usuário", err)
	}

	return &RegisterUserOutput{
		Message: "Usuário registrado com sucesso",
		Detail:  "O usuário foi criado e registrado no sistema com sucesso",
	}, nil
}
