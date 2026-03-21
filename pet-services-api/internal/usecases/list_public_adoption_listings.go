package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListPublicAdoptionListingsInput struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Sex      string `json:"sex" form:"sex"`
	Size     string `json:"size" form:"size"`
	AgeGroup string `json:"age_group" form:"age_group"`
	CityID   string `json:"city_id" form:"city_id"`
	StateID  string `json:"state_id" form:"state_id"`
}

type ListPublicAdoptionListingsOutput struct {
	Listings   []*entities.AdoptionListing `json:"listings"`
	Pagination PaginationMetadata          `json:"pagination"`
}

type ListPublicAdoptionListingsUseCase struct {
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewListPublicAdoptionListingsUseCase(
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *ListPublicAdoptionListingsUseCase {
	return &ListPublicAdoptionListingsUseCase{
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *ListPublicAdoptionListingsUseCase) Execute(ctx context.Context, input ListPublicAdoptionListingsInput) (*ListPublicAdoptionListingsOutput, []exceptions.ProblemDetails) {
	const from = "ListPublicAdoptionListingsUseCase.Execute"

	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 12
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}

	filters := entities.AdoptionListingFilters{
		Sex:      input.Sex,
		Size:     input.Size,
		AgeGroup: input.AgeGroup,
		CityID:   input.CityID,
		StateID:  input.StateID,
	}

	listings, total, err := uc.adoptionListingRepository.ListPublic(filters, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar anúncios de adoção", err)
	}

	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize != 0 {
		totalPages++
	}

	return &ListPublicAdoptionListingsOutput{
		Listings: listings,
		Pagination: PaginationMetadata{
			CurrentPage:  input.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      input.PageSize,
		},
	}, nil
}
