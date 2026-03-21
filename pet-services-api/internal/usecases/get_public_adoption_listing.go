package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type GetPublicAdoptionListingInput struct {
	ListingID string `json:"listing_id"`
}

type GetPublicAdoptionListingOutput struct {
	Listing *entities.AdoptionListing `json:"listing"`
}

type GetPublicAdoptionListingUseCase struct {
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewGetPublicAdoptionListingUseCase(
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *GetPublicAdoptionListingUseCase {
	return &GetPublicAdoptionListingUseCase{
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *GetPublicAdoptionListingUseCase) Execute(ctx context.Context, input GetPublicAdoptionListingInput) (*GetPublicAdoptionListingOutput, []exceptions.ProblemDetails) {
	const from = "GetPublicAdoptionListingUseCase.Execute"

	if input.ListingID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do anúncio ausente", errors.New("O ID do anúncio é obrigatório"))
	}

	listing, err := uc.adoptionListingRepository.FindByID(input.ListingID)
	if err != nil {
		if err.Error() == consts.AdoptionListingNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Anúncio não encontrado", errors.New("O anúncio de adoção informado não foi encontrado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar anúncio", err)
	}

	if listing.Status != entities.AdoptionListingStatuses.Published {
		return nil, uc.logger.LogNotFound(ctx, from, "Anúncio não disponível", errors.New("Este anúncio não está disponível publicamente"))
	}

	return &GetPublicAdoptionListingOutput{Listing: listing}, nil
}
