package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type WithdrawAdoptionApplicationInput struct {
	ApplicationID string
	UserID        string
}

type WithdrawAdoptionApplicationOutput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type WithdrawAdoptionApplicationUseCase struct {
	applicationRepo entities.AdoptionApplicationRepository
	logger          logging.LoggerInterface
}

func NewWithdrawAdoptionApplicationUseCase(
	applicationRepo entities.AdoptionApplicationRepository,
	logger logging.LoggerInterface,
) *WithdrawAdoptionApplicationUseCase {
	return &WithdrawAdoptionApplicationUseCase{
		applicationRepo: applicationRepo,
		logger:          logger,
	}
}

func (u *WithdrawAdoptionApplicationUseCase) Execute(ctx context.Context, input WithdrawAdoptionApplicationInput) (*WithdrawAdoptionApplicationOutput, []exceptions.ProblemDetails) {
	// Buscar candidatura
	application, err := u.applicationRepo.FindByID(input.ApplicationID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionApplicationNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "WithdrawAdoptionApplicationUseCase", "Candidatura não encontrada", errors.New("A candidatura especificada não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar candidatura",
			Detail: "Ocorreu um erro ao buscar a candidatura",
		})
		u.logger.LogError(ctx, "WithdrawAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Verificar permissão (apenas o candidato pode desistir)
	if application.ApplicantUserID != input.UserID {
		return nil, u.logger.LogForbidden(ctx, "WithdrawAdoptionApplicationUseCase", "Permissão negada", errors.New("Apenas o candidato pode desistir desta candidatura"))
	}

	// Verificar se já foi rejeitada, aprovada, retirada
	if application.Status == entities.AdoptionApplicationStatuses.Withdrawn ||
		application.Status == entities.AdoptionApplicationStatuses.Approved ||
		application.Status == entities.AdoptionApplicationStatuses.Rejected {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Candidatura não pode ser retirada",
			Detail: "Esta candidatura não pode mais ser retirada",
		})
		u.logger.LogBadRequest(ctx, "WithdrawAdoptionApplicationUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Retirar candidatura
	application.Withdraw()

	// Persistir
	if err := u.applicationRepo.Update(application); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar candidatura",
			Detail: "Ocorreu um erro ao retirar a candidatura",
		})
		u.logger.LogError(ctx, "WithdrawAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "WithdrawAdoptionApplicationUseCase", "Candidatura "+input.ApplicationID+" retirada")

	return &WithdrawAdoptionApplicationOutput{
		ID:     application.ID,
		Status: application.Status,
	}, nil
}
