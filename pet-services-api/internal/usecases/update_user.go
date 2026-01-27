package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type UpdateUserInput struct {
	UserID   string           `json:"user_id"`
	Name     string           `json:"name,omitempty"`
	UserType string           `json:"user_type,omitempty"`
	Phone    entities.Phone   `json:"phone,omitempty"`
	Address  entities.Address `json:"address,omitempty"`
}

type UpdateUserOutput struct {
	Message string         `json:"message,omitempty"`
	Detail  string         `json:"detail,omitempty"`
	User    *entities.User `json:"user,omitempty"`
}

type UpdateUserUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewUpdateUserUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, []exceptions.ProblemDetails) {
	const from = "UpdateUserUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", nil)
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.UserType != "" {
		user.UserType = input.UserType
	}

	if input.Phone.CountryCode != "" || input.Phone.AreaCode != "" || input.Phone.Number != "" {
		phone, problems := entities.NewPhone(input.Phone.CountryCode, input.Phone.AreaCode, input.Phone.Number)
		if len(problems) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Telefone inválido", problems)
			return nil, problems
		}
		user.Phone = *phone
	}

	if input.Address.Street != "" || input.Address.Number != "" || input.Address.Neighborhood != "" ||
		input.Address.City != "" || input.Address.ZipCode != "" || input.Address.State != "" ||
		input.Address.Country != "" || input.Address.Complement != "" ||
		input.Address.Location.Latitude != 0 || input.Address.Location.Longitude != 0 {

		loc, problems := entities.NewLocation(input.Address.Location.Latitude, input.Address.Location.Longitude)
		if len(problems) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Localização inválida", problems)
			return nil, problems
		}

		addr, problems := entities.NewAddress(
			input.Address.Street,
			input.Address.Number,
			input.Address.Neighborhood,
			input.Address.City,
			input.Address.ZipCode,
			input.Address.State,
			input.Address.Country,
			input.Address.Complement,
			*loc,
		)
		if len(problems) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Endereço inválido", problems)
			return nil, problems
		}
		user.Address = *addr
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar usuário", err)
	}

	return &UpdateUserOutput{
		Message: "Usuário atualizado com sucesso",
		Detail:  "Os dados do usuário foram atualizados",
		User:    user,
	}, nil
}
