package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type UpdateUserInputBody struct {
	Name     string           `json:"name,omitempty"`
	UserType string           `json:"user_type,omitempty"`
	Phone    entities.Phone   `json:"phone"`
	Address  entities.Address `json:"address"`
}

type UpdateUserInput struct {
	UserID   string           `json:"user_id"`
	Name     string           `json:"name,omitempty"`
	UserType string           `json:"user_type,omitempty"`
	Phone    entities.Phone   `json:"phone"`
	Address  entities.Address `json:"address"`
}

type UpdateUserOutput struct {
	Message string      `json:"message,omitempty"`
	Detail  string      `json:"detail,omitempty"`
	User    *UserOutput `json:"user,omitempty"`
}

type UpdateUserUseCase struct {
	userRepository entities.UserRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewUpdateUserUseCase(userRepo entities.UserRepository, storageService storage.ObjectStorage, logger logging.LoggerInterface) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepo,
		storage:        storageService,
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
		if input.UserType == "admin" {
			return nil, uc.logger.LogForbidden(ctx, from, "Tipo de usuário não permitido", errors.New("Não é permitido definir usuário como admin"))
		}

		if input.UserType != entities.UserTypes.Owner && input.UserType != entities.UserTypes.Provider {
			return nil, uc.logger.LogBadRequest(ctx, from, "Tipo de usuário inválido", errors.New("O tipo de usuário deve ser 'owner' ou 'provider'"))
		}

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

	user.ProfileComplete = user.UserType != "" &&
		user.Name != "" &&
		user.Phone.CountryCode != "" &&
		user.Phone.AreaCode != "" &&
		user.Phone.Number != "" &&
		user.Address.Street != "" &&
		user.Address.Number != "" &&
		user.Address.Neighborhood != "" &&
		user.Address.City != "" &&
		user.Address.ZipCode != "" &&
		user.Address.State != "" &&
		user.Address.Country != ""

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar usuário", err)
	}

	if err := storage.SignUserPhotos(ctx, uc.storage, user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URLs das fotos", err)
	}

	return &UpdateUserOutput{
		Message: "Usuário atualizado com sucesso",
		Detail:  "Os dados do usuário foram atualizados",
		User:    NewUserOutput(user),
	}, nil
}
