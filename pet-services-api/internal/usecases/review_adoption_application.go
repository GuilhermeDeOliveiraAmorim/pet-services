package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ReviewAdoptionApplicationInputBody struct {
	Action        string `json:"action" binding:"required"` // under_review, interview, approve, reject
	NotesInternal string `json:"notes_internal,omitempty"`
}

type ReviewAdoptionApplicationInput struct {
	ApplicationID string
	ReviewerID    string
	ReviewAdoptionApplicationInputBody
}

type ReviewAdoptionApplicationOutput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ReviewAdoptionApplicationUseCase struct {
	applicationRepo entities.AdoptionApplicationRepository
	logger          logging.LoggerInterface
}

func NewReviewAdoptionApplicationUseCase(
	applicationRepo entities.AdoptionApplicationRepository,
	logger logging.LoggerInterface,
) *ReviewAdoptionApplicationUseCase {
	return &ReviewAdoptionApplicationUseCase{
		applicationRepo: applicationRepo,
		logger:          logger,
	}
}

func (u *ReviewAdoptionApplicationUseCase) Execute(ctx context.Context, input ReviewAdoptionApplicationInput) (*ReviewAdoptionApplicationOutput, []exceptions.ProblemDetails) {
	// Buscar candidatura
	application, err := u.applicationRepo.FindByID(input.ApplicationID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionApplicationNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "ReviewAdoptionApplicationUseCase", "Candidatura não encontrada", errors.New("A candidatura especificada não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar candidatura",
			Detail: "Ocorreu um erro ao buscar a candidatura",
		})
		u.logger.LogError(ctx, "ReviewAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Validar ação
	if input.Action == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Ação ausente",
			Detail: "A ação é obrigatória",
		})
		u.logger.LogBadRequest(ctx, "ReviewAdoptionApplicationUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Aplicar ação
	switch input.Action {
	case "under_review":
		application.MoveToUnderReview(input.ReviewerID)
	case "interview":
		application.MoveToInterview(input.ReviewerID)
	case "approve":
		application.Approve(input.ReviewerID)
	case "reject":
		application.Reject(input.ReviewerID, input.NotesInternal)
	default:
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Ação inválida",
			Detail: "A ação especificada não é válida",
		})
		u.logger.LogBadRequest(ctx, "ReviewAdoptionApplicationUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Persistir
	if err := u.applicationRepo.Update(application); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar candidatura",
			Detail: "Ocorreu um erro ao atualizar a candidatura",
		})
		u.logger.LogError(ctx, "ReviewAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "ReviewAdoptionApplicationUseCase", "Candidatura "+input.ApplicationID+" movida para "+application.Status)

	return &ReviewAdoptionApplicationOutput{
		ID:     application.ID,
		Status: application.Status,
	}, nil
}
