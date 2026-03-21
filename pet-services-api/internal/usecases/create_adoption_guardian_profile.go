package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CreateAdoptionGuardianProfileInputBody struct {
	DisplayName  string `json:"display_name"`
	GuardianType string `json:"guardian_type"`
	Document     string `json:"document"`
	Phone        string `json:"phone"`
	Whatsapp     string `json:"whatsapp"`
	About        string `json:"about"`
	CityID       string `json:"city_id"`
	StateID      string `json:"state_id"`
}

type CreateAdoptionGuardianProfileInput struct {
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

type CreateAdoptionGuardianProfileOutput struct {
	Message string                           `json:"message,omitempty"`
	Detail  string                           `json:"detail,omitempty"`
	Profile *entities.AdoptionGuardianProfile `json:"profile,omitempty"`
}

type CreateAdoptionGuardianProfileUseCase struct {
	userRepository                   entities.UserRepository
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository
	logger                           logging.LoggerInterface
}

func NewCreateAdoptionGuardianProfileUseCase(
	userRepository entities.UserRepository,
	adoptionGuardianProfileRepository entities.AdoptionGuardianProfileRepository,
	logger logging.LoggerInterface,
) *CreateAdoptionGuardianProfileUseCase {
	return &CreateAdoptionGuardianProfileUseCase{
		userRepository:                   userRepository,
		adoptionGuardianProfileRepository: adoptionGuardianProfileRepository,
		logger:                           logger,
	}
}

func (uc *CreateAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input CreateAdoptionGuardianProfileInput) (*CreateAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
	const from = "CreateAdoptionGuardianProfileUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	_, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	existing, err := uc.adoptionGuardianProfileRepository.FindByUserID(input.UserID)
	if err == nil && existing != nil {
		return nil, uc.logger.LogConflict(ctx, from, "Perfil já existe", errors.New("O usuário já possui um perfil de responsável por adoção cadastrado"))
	}
	if err != nil && err.Error() != consts.AdoptionGuardianProfileNotFoundError {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar perfil existente", err)
	}

	profile, problems := entities.NewAdoptionGuardianProfile(
		input.UserID,
		input.DisplayName,
		input.GuardianType,
		input.Document,
		input.Phone,
		input.Whatsapp,
		input.About,
		input.CityID,
		input.StateID,
	)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Dados de perfil inválidos", problems)
		return nil, problems
	}

	if err := uc.adoptionGuardianProfileRepository.Create(profile); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar perfil de responsável", err)
	}

	return &CreateAdoptionGuardianProfileOutput{
		Message: "Perfil de responsável criado com sucesso",
		Detail:  "O perfil foi criado e está aguardando aprovação",
		Profile: profile,
	}, nil
}
