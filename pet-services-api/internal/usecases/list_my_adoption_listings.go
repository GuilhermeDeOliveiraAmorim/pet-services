package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListMyAdoptionListingsInput struct {
	GuardianProfileID string `json:"guardian_profile_id"`
	Page              int    `json:"page" form:"page"`
	PageSize          int    `json:"page_size" form:"page_size"`
}

type ListMyAdoptionListingsOutput struct {
	Listings   []*entities.AdoptionListing `json:"listings"`
	Pagination PaginationMetadata          `json:"pagination"`
}

type ListMyAdoptionListingsUseCase struct {
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewListMyAdoptionListingsUseCase(
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *ListMyAdoptionListingsUseCase {
	return &ListMyAdoptionListingsUseCase{
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *ListMyAdoptionListingsUseCase) Execute(ctx context.Context, input ListMyAdoptionListingsInput) (*ListMyAdoptionListingsOutput, []exceptions.ProblemDetails) {
	const from = "ListMyAdoptionListingsUseCase.Execute"

	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}

	listings, total, err := uc.adoptionListingRepository.ListByGuardianProfileID(input.GuardianProfileID, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar meus anúncios de adoção", err)
	}

	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize != 0 {
		totalPages++
	}

	return &ListMyAdoptionListingsOutput{
		Listings: listings,
		Pagination: PaginationMetadata{
			CurrentPage:  input.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      input.PageSize,
		},
	}, nil
}
