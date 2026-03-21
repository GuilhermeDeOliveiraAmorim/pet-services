package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListAdoptionApplicationsByListingInput struct {
	ListingID string
	Page      int
	PageSize  int
}

type ListAdoptionApplicationsByListingOutput struct {
	Applications []*AdoptionApplicationDTO `json:"applications"`
	Pagination   PaginationMetadata        `json:"pagination"`
}

type ListAdoptionApplicationsByListingUseCase struct {
	listingRepository entities.AdoptionListingRepository
	applicationRepo   entities.AdoptionApplicationRepository
	logger            logging.LoggerInterface
}

func NewListAdoptionApplicationsByListingUseCase(
	listingRepository entities.AdoptionListingRepository,
	applicationRepo entities.AdoptionApplicationRepository,
	logger logging.LoggerInterface,
) *ListAdoptionApplicationsByListingUseCase {
	return &ListAdoptionApplicationsByListingUseCase{
		listingRepository: listingRepository,
		applicationRepo:   applicationRepo,
		logger:            logger,
	}
}

func (u *ListAdoptionApplicationsByListingUseCase) Execute(ctx context.Context, input ListAdoptionApplicationsByListingInput) (*ListAdoptionApplicationsByListingOutput, []exceptions.ProblemDetails) {
	_, err := u.listingRepository.FindByID(input.ListingID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionListingNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "ListAdoptionApplicationsByListingUseCase", "Anúncio não encontrado", errors.New("O anúncio especificado não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar anúncio",
			Detail: "Ocorreu um erro ao buscar o anúncio",
		})
		u.logger.LogError(ctx, "ListAdoptionApplicationsByListingUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}

	applications, total, err := u.applicationRepo.ListByListingID(input.ListingID, input.Page, input.PageSize)
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao listar candidaturas",
			Detail: "Ocorreu um erro ao listar as candidaturas do anúncio",
		})
		u.logger.LogError(ctx, "ListAdoptionApplicationsByListingUseCase", problem.Detail, err)
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

	return &ListAdoptionApplicationsByListingOutput{
		Applications: dtos,
		Pagination: PaginationMetadata{
			CurrentPage:  input.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      input.PageSize,
		},
	}, nil
}
