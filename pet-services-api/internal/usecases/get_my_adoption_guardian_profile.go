package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type GetMyAdoptionGuardianProfileInput struct {
	UserID string `json:"user_id"`
}

type GetMyAdoptionGuardianProfileOutput struct {
	Profile *entities.AdoptionGuardianProfile `json:"profile"`
}

type GetMyAdoptionGuardianProfileUseCase struct {
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository
	logger                            logging.LoggerInterface
}

func NewGetMyAdoptionGuardianProfileUseCase(
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository,
	logger logging.LoggerInterface,
) *GetMyAdoptionGuardianProfileUseCase {
	return &GetMyAdoptionGuardianProfileUseCase{
		adoptionGuardianProfileRepository: adoptionGuardianProfileRepository,
		logger:                            logger,
	}
}

func (uc *GetMyAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input GetMyAdoptionGuardianProfileInput) (*GetMyAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
	const from = "GetMyAdoptionGuardianProfileUseCase.Execute"

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

	return &GetMyAdoptionGuardianProfileOutput{Profile: profile}, nil
}
