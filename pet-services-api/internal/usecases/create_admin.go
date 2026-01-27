package usecases

import (
	"context"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CreateAdminInput struct {
	Name    string           `json:"name"`
	Login   entities.Login   `json:"login"`
	Phone   entities.Phone   `json:"phone"`
	Address entities.Address `json:"address"`
}

type CreateAdminOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type CreateAdminUseCase struct {
	userRepository entities.UserRepository
}

func NewCreateAdminUseCase(userRepository entities.UserRepository) *CreateAdminUseCase {
	return &CreateAdminUseCase{
		userRepository: userRepository,
	}
}

func (uc *CreateAdminUseCase) Execute(ctx context.Context, input CreateAdminInput, isAdmin bool) (*CreateAdminOutput, []exceptions.ProblemDetails) {
	const from = "CreateAdminUseCase.Execute"

	if !isAdmin {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Apenas administradores podem criar outros administradores.",
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

	problems := []exceptions.ProblemDetails{}

	login, problems := entities.NewLogin(input.Login.Email, input.Login.Password)
	if len(problems) > 0 {
		problems = append(problems, problems...)
	}

	if err := login.EncryptPassword(); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	phone, problems := entities.NewPhone(input.Phone.CountryCode, input.Phone.AreaCode, input.Phone.Number)
	if len(problems) > 0 {
		problems = append(problems, problems...)
	}

	location, problems := entities.NewLocation(input.Address.Location.Latitude, input.Address.Location.Longitude)
	if len(problems) > 0 {
		problems = append(problems, problems...)
	}

	address, problems := entities.NewAddress(
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

	if len(problems) > 0 {
		problems = append(problems, problems...)
	}

	user, problems := entities.NewUser(
		input.Name,
		"admin",
		*login,
		*phone,
		*address,
	)

	if len(problems) > 0 {
		problems = append(problems, problems...)
	}

	if len(problems) > 0 {
		return nil, problems
	}

	if err := uc.userRepository.Create(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criar admin", err)
	}

	return &CreateAdminOutput{
		Message: "Administrador criado com sucesso",
		Detail:  "O usuário administrador foi criado no sistema com sucesso",
	}, nil
}
