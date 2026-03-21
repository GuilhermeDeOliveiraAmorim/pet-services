package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListMyAdoptionApplicationsInput struct {
	UserID   string
	Page     int
	PageSize int
}

type AdoptionApplicationDTO struct {
	ID           string `json:"id"`
	ListingID    string `json:"listing_id"`
	Status       string `json:"status"`
	Motivation   string `json:"motivation"`
	ContactPhone string `json:"contact_phone"`
	ReviewedAt   string `json:"reviewed_at,omitempty"`
	CreatedAt    string `json:"created_at"`
}

type ListMyAdoptionApplicationsOutput struct {
	Applications []*AdoptionApplicationDTO `json:"applications"`
	Pagination   PaginationMetadata        `json:"pagination"`
}

type ListMyAdoptionApplicationsUseCase struct {
	applicationRepo entities.AdoptionApplicationRepository
	logger          logging.LoggerInterface
}

func NewListMyAdoptionApplicationsUseCase(
	applicationRepo entities.AdoptionApplicationRepository,
	logger logging.LoggerInterface,
) *ListMyAdoptionApplicationsUseCase {
	return &ListMyAdoptionApplicationsUseCase{
		applicationRepo: applicationRepo,
		logger:          logger,
	}
}

func (u *ListMyAdoptionApplicationsUseCase) Execute(ctx context.Context, input ListMyAdoptionApplicationsInput) (*ListMyAdoptionApplicationsOutput, []exceptions.ProblemDetails) {
	// Validar paginação
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}

	applications, total, err := u.applicationRepo.ListByApplicantUserID(input.UserID, input.Page, input.PageSize)
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao listar candidaturas",
			Detail: "Ocorreu um erro ao listar suas candidaturas",
		})
		u.logger.LogError(ctx, "ListMyAdoptionApplicationsUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	dtos := make([]*AdoptionApplicationDTO, len(applications))
	for i, app := range applications {
		reviewedAt := ""
		if app.ReviewedAt != nil {
			reviewedAt = app.ReviewedAt.Format("2006-01-02T15:04:05Z")
		}
		dtos[i] = &AdoptionApplicationDTO{
			ID:           app.ID,
			ListingID:    app.ListingID,
			Status:       app.Status,
			Motivation:   app.Motivation,
			ContactPhone: app.ContactPhone,
			ReviewedAt:   reviewedAt,
			CreatedAt:    app.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	totalPages := (int(total) + input.PageSize - 1) / input.PageSize

	return &ListMyAdoptionApplicationsOutput{
		Applications: dtos,
		Pagination: PaginationMetadata{
			CurrentPage:  input.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      input.PageSize,
		},
	}, nil
}
