package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ChangeAdoptionListingStatusInput struct {
	ListingID         string `json:"listing_id"`
	GuardianProfileID string `json:"guardian_profile_id"`
	Action            string `json:"action"` // publish | pause
}

type ChangeAdoptionListingStatusOutput struct {
	Message string                    `json:"message"`
	Listing *entities.AdoptionListing `json:"listing"`
}

type ChangeAdoptionListingStatusUseCase struct {
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewChangeAdoptionListingStatusUseCase(
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *ChangeAdoptionListingStatusUseCase {
	return &ChangeAdoptionListingStatusUseCase{
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *ChangeAdoptionListingStatusUseCase) Execute(ctx context.Context, input ChangeAdoptionListingStatusInput) (*ChangeAdoptionListingStatusOutput, []exceptions.ProblemDetails) {
	const from = "ChangeAdoptionListingStatusUseCase.Execute"

	listing, err := uc.adoptionListingRepository.FindByID(input.ListingID)
	if err != nil {
		if err.Error() == consts.AdoptionListingNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Anúncio não encontrado", errors.New("O anúncio de adoção informado não foi encontrado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar anúncio", err)
	}

	if listing.GuardianProfileID != input.GuardianProfileID {
		return nil, uc.logger.LogForbidden(ctx, from, "Sem permissão", errors.New("Você não tem permissão para alterar o status deste anúncio"))
	}

	var message string
	switch input.Action {
	case "publish":
		if listing.Status != entities.AdoptionListingStatuses.Draft && listing.Status != entities.AdoptionListingStatuses.Paused {
			return nil, uc.logger.LogBadRequest(ctx, from, "Transição inválida", errors.New("Somente anúncios com status 'draft' ou 'paused' podem ser publicados"))
		}
		listing.Publish()
		message = "Anúncio publicado com sucesso"
	case "pause":
		if listing.Status != entities.AdoptionListingStatuses.Published {
			return nil, uc.logger.LogBadRequest(ctx, from, "Transição inválida", errors.New("Somente anúncios publicados podem ser pausados"))
		}
		listing.Pause()
		message = "Anúncio pausado com sucesso"
	case "archive":
		if listing.Status == entities.AdoptionListingStatuses.Adopted {
			return nil, uc.logger.LogBadRequest(ctx, from, "Transição inválida", errors.New("Anúncios já adotados não podem ser arquivados diretamente"))
		}
		listing.Archive()
		message = "Anúncio arquivado com sucesso"
	default:
		return nil, uc.logger.LogBadRequest(ctx, from, "Ação inválida", errors.New("A ação informada não é válida. Use: publish, pause ou archive"))
	}

	if err := uc.adoptionListingRepository.Update(listing); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar status do anúncio", err)
	}

	return &ChangeAdoptionListingStatusOutput{
		Message: message,
		Listing: listing,
	}, nil
}
