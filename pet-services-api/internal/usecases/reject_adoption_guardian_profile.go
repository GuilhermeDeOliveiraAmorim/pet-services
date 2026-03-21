package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type RejectAdoptionGuardianProfileInputBody struct {
	Reason string `json:"reason,omitempty"`
}

type RejectAdoptionGuardianProfileInput struct {
	ProfileID  string
	ReviewedBy string
	RejectAdoptionGuardianProfileInputBody
}

type RejectAdoptionGuardianProfileOutput struct {
	ID             string `json:"id"`
	ApprovalStatus string `json:"approval_status"`
}

type RejectAdoptionGuardianProfileUseCase struct {
	guardianProfileRepo entities.AdoptionGuardianProfileRepository
	logger              logging.LoggerInterface
}

func NewRejectAdoptionGuardianProfileUseCase(
	guardianProfileRepo entities.AdoptionGuardianProfileRepository,
	logger logging.LoggerInterface,
) *RejectAdoptionGuardianProfileUseCase {
	return &RejectAdoptionGuardianProfileUseCase{
		guardianProfileRepo: guardianProfileRepo,
		logger:              logger,
	}
}

func (u *RejectAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input RejectAdoptionGuardianProfileInput) (*RejectAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
	// Buscar perfil
	profile, err := u.guardianProfileRepo.FindByID(input.ProfileID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionGuardianProfileNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "RejectAdoptionGuardianProfileUseCase", "Perfil não encontrado", errors.New("O perfil especificado não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar perfil",
			Detail: "Ocorreu um erro ao buscar o perfil",
		})
		u.logger.LogError(ctx, "RejectAdoptionGuardianProfileUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Verificar se está pendente
	if profile.ApprovalStatus != entities.AdoptionGuardianApprovalStatuses.Pending {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Perfil já foi analisado",
			Detail: "Este perfil já foi aprovado ou rejeitado",
		})
		u.logger.LogBadRequest(ctx, "RejectAdoptionGuardianProfileUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Rejeitar
	profile.Reject(input.ReviewedBy)

	// Persistir
	if err := u.guardianProfileRepo.Update(profile); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar perfil",
			Detail: "Ocorreu um erro ao rejeitar o perfil",
		})
		u.logger.LogError(ctx, "RejectAdoptionGuardianProfileUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "RejectAdoptionGuardianProfileUseCase", "Perfil "+input.ProfileID+" rejeitado")

	return &RejectAdoptionGuardianProfileOutput{
		ID:             profile.ID,
		ApprovalStatus: profile.ApprovalStatus,
	}, nil
}
