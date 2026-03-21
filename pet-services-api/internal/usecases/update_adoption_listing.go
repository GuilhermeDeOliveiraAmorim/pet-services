package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type UpdateAdoptionListingInputBody struct {
	Title                  string  `json:"title"`
	Description            string  `json:"description"`
	AdoptionReason         string  `json:"adoption_reason"`
	Sex                    string  `json:"sex"`
	Size                   string  `json:"size"`
	AgeGroup               string  `json:"age_group"`
	CityID                 string  `json:"city_id"`
	StateID                string  `json:"state_id"`
	Latitude               float64 `json:"latitude"`
	Longitude              float64 `json:"longitude"`
	Vaccinated             *bool   `json:"vaccinated"`
	Neutered               *bool   `json:"neutered"`
	Dewormed               *bool   `json:"dewormed"`
	SpecialNeeds           *bool   `json:"special_needs"`
	GoodWithChildren       *bool   `json:"good_with_children"`
	GoodWithDogs           *bool   `json:"good_with_dogs"`
	GoodWithCats           *bool   `json:"good_with_cats"`
	RequiresHouseScreening *bool   `json:"requires_house_screening"`
}

type UpdateAdoptionListingInput struct {
	ListingID         string `json:"listing_id"`
	GuardianProfileID string `json:"guardian_profile_id"`
	UpdateAdoptionListingInputBody
}

type UpdateAdoptionListingOutput struct {
	Message string                    `json:"message"`
	Listing *entities.AdoptionListing `json:"listing"`
}

type UpdateAdoptionListingUseCase struct {
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewUpdateAdoptionListingUseCase(
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *UpdateAdoptionListingUseCase {
	return &UpdateAdoptionListingUseCase{
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *UpdateAdoptionListingUseCase) Execute(ctx context.Context, input UpdateAdoptionListingInput) (*UpdateAdoptionListingOutput, []exceptions.ProblemDetails) {
	const from = "UpdateAdoptionListingUseCase.Execute"

	listing, err := uc.adoptionListingRepository.FindByID(input.ListingID)
	if err != nil {
		if err.Error() == consts.AdoptionListingNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Anúncio não encontrado", errors.New("O anúncio de adoção informado não foi encontrado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar anúncio", err)
	}

	if listing.GuardianProfileID != input.GuardianProfileID {
		return nil, uc.logger.LogForbidden(ctx, from, "Sem permissão", errors.New("Você não tem permissão para editar este anúncio"))
	}

	if listing.Status != entities.AdoptionListingStatuses.Draft && listing.Status != entities.AdoptionListingStatuses.Paused {
		return nil, uc.logger.LogBadRequest(ctx, from, "Anúncio não pode ser editado", errors.New("Somente anúncios com status 'draft' ou 'paused' podem ser editados"))
	}

	if input.Title != "" {
		listing.Title = input.Title
	}
	if input.Description != "" {
		listing.Description = input.Description
	}
	if input.AdoptionReason != "" {
		listing.AdoptionReason = input.AdoptionReason
	}
	if input.Sex != "" {
		listing.Sex = input.Sex
	}
	if input.Size != "" {
		listing.Size = input.Size
	}
	if input.AgeGroup != "" {
		listing.AgeGroup = input.AgeGroup
	}
	if input.CityID != "" {
		listing.CityID = input.CityID
	}
	if input.StateID != "" {
		listing.StateID = input.StateID
	}
	if input.Latitude != 0 {
		listing.Latitude = input.Latitude
	}
	if input.Longitude != 0 {
		listing.Longitude = input.Longitude
	}
	if input.Vaccinated != nil {
		listing.Vaccinated = *input.Vaccinated
	}
	if input.Neutered != nil {
		listing.Neutered = *input.Neutered
	}
	if input.Dewormed != nil {
		listing.Dewormed = *input.Dewormed
	}
	if input.SpecialNeeds != nil {
		listing.SpecialNeeds = *input.SpecialNeeds
	}
	if input.GoodWithChildren != nil {
		listing.GoodWithChildren = *input.GoodWithChildren
	}
	if input.GoodWithDogs != nil {
		listing.GoodWithDogs = *input.GoodWithDogs
	}
	if input.GoodWithCats != nil {
		listing.GoodWithCats = *input.GoodWithCats
	}
	if input.RequiresHouseScreening != nil {
		listing.RequiresHouseScreening = *input.RequiresHouseScreening
	}

	if err := uc.adoptionListingRepository.Update(listing); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar anúncio", err)
	}

	return &UpdateAdoptionListingOutput{
		Message: "Anúncio de adoção atualizado com sucesso",
		Listing: listing,
	}, nil
}
