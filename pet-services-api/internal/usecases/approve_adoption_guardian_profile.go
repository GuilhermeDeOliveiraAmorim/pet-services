package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
)

type ApproveAdoptionGuardianProfileInput struct {
	ProfileID  string
	ApprovedBy string
}

type ApproveAdoptionGuardianProfileOutput struct {
	ID             string `json:"id"`
	ApprovalStatus string `json:"approval_status"`
}

type ApproveAdoptionGuardianProfileUseCase struct {
	guardianProfileRepo entities.AdoptionGuardianProfileRepository
	userRepository      entities.UserRepository
	emailService        mail.EmailService
	logger              logging.LoggerInterface
}

func NewApproveAdoptionGuardianProfileUseCase(
	guardianProfileRepo entities.AdoptionGuardianProfileRepository,
	userRepository entities.UserRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *ApproveAdoptionGuardianProfileUseCase {
	return &ApproveAdoptionGuardianProfileUseCase{
		guardianProfileRepo: guardianProfileRepo,
		userRepository:      userRepository,
		emailService:        emailService,
		logger:              logger,
	}
}

func (u *ApproveAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input ApproveAdoptionGuardianProfileInput) (*ApproveAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
	profile, err := u.guardianProfileRepo.FindByID(input.ProfileID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionGuardianProfileNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "ApproveAdoptionGuardianProfileUseCase", "Perfil não encontrado", errors.New("O perfil especificado não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar perfil",
			Detail: "Ocorreu um erro ao buscar o perfil",
		})
		u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	if profile.ApprovalStatus != entities.AdoptionGuardianApprovalStatuses.Pending {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Perfil já foi analisado",
			Detail: "Este perfil já foi aprovado ou rejeitado",
		})
		u.logger.LogBadRequest(ctx, "ApproveAdoptionGuardianProfileUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	profile.Approve(input.ApprovedBy)

	if err := u.guardianProfileRepo.Update(profile); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar perfil",
			Detail: "Ocorreu um erro ao aprovar o perfil",
		})
		u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "ApproveAdoptionGuardianProfileUseCase", "Perfil "+input.ProfileID+" aprovado")

	user, err := u.userRepository.FindByID(profile.UserID)
	if err == nil {
		if err := u.emailService.SendAdoptionGuardianProfileApprovedEmail(user.Login.Email, user.Name); err != nil {
			u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", "Erro ao enviar email de aprovação", err)

		}
	} else {
		u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", "Erro ao buscar usuário para enviar email", err)
	}

	return &ApproveAdoptionGuardianProfileOutput{
		ID:             profile.ID,
		ApprovalStatus: profile.ApprovalStatus,
	}, nil
}
