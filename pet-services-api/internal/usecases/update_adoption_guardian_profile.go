package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type UpdateAdoptionGuardianProfileInputBody struct {
	DisplayName  string `json:"display_name"`
	GuardianType string `json:"guardian_type"`
	Document     string `json:"document"`
	Phone        string `json:"phone"`
	Whatsapp     string `json:"whatsapp"`
	About        string `json:"about"`
	CityID       string `json:"city_id"`
	StateID      string `json:"state_id"`
}

type UpdateAdoptionGuardianProfileInput struct {
	UserID       string `json:"user_id"`
	DisplayName  string `json:"display_name"`
	GuardianType string `json:"guardian_type"`
	Document     string `json:"document"`
	Phone        string `json:"phone"`
	Whatsapp     string `json:"whatsapp"`
	About        string `json:"about"`
	CityID       string `json:"city_id"`
	StateID      string `json:"state_id"`
}

type UpdateAdoptionGuardianProfileOutput struct {
	Message string                            `json:"message"`
	Profile *entities.AdoptionGuardianProfile `json:"profile"`
}

type UpdateAdoptionGuardianProfileUseCase struct {
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository
	logger                            logging.LoggerInterface
}

func NewUpdateAdoptionGuardianProfileUseCase(
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository,
	logger logging.LoggerInterface,
) *UpdateAdoptionGuardianProfileUseCase {
	return &UpdateAdoptionGuardianProfileUseCase{
		adoptionGuardianProfileRepository: adoptionGuardianProfileRepository,
		logger:                            logger,
	}
}

func (uc *UpdateAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input UpdateAdoptionGuardianProfileInput) (*UpdateAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
	const from = "UpdateAdoptionGuardianProfileUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	profile, err := uc.adoptionGuardianProfileRepository.FindByUserID(input.UserID)
	if err != nil {
		if err.Error() == consts.AdoptionGuardianProfileNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Perfil não encontrado", errors.New("Nenhum perfil de responsável por adoção encontrado para este usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar perfil", err)
	}

	if input.DisplayName != "" {
		profile.DisplayName = input.DisplayName
	}
	if input.GuardianType != "" {
		gt := input.GuardianType
		if gt != entities.AdoptionGuardianTypes.NGO && gt != entities.AdoptionGuardianTypes.Independent && gt != entities.AdoptionGuardianTypes.Owner {
			return nil, uc.logger.LogBadRequest(ctx, from, "Tipo de responsável inválido", errors.New("O tipo de responsável informado não é válido. Use: ngo, independent ou owner"))
		}
		profile.GuardianType = gt
	}
	if input.Document != "" {
		profile.Document = input.Document
	}
	if input.Phone != "" {
		profile.Phone = input.Phone
	}
	if input.Whatsapp != "" {
		profile.Whatsapp = input.Whatsapp
	}
	if input.About != "" {
		profile.About = input.About
	}
	if input.CityID != "" {
		profile.CityID = input.CityID
	}
	if input.StateID != "" {
		profile.StateID = input.StateID
	}

	if err := uc.adoptionGuardianProfileRepository.Update(profile); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar perfil", err)
	}

	return &UpdateAdoptionGuardianProfileOutput{
		Message: "Perfil de responsável atualizado com sucesso",
		Profile: profile,
	}, nil
}
