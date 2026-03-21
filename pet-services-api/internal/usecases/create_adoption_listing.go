package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CreateAdoptionListingInputBody struct {
	PetID                  string  `json:"pet_id"`
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
	Vaccinated             bool    `json:"vaccinated"`
	Neutered               bool    `json:"neutered"`
	Dewormed               bool    `json:"dewormed"`
	SpecialNeeds           bool    `json:"special_needs"`
	GoodWithChildren       bool    `json:"good_with_children"`
	GoodWithDogs           bool    `json:"good_with_dogs"`
	GoodWithCats           bool    `json:"good_with_cats"`
	RequiresHouseScreening bool    `json:"requires_house_screening"`
}

type CreateAdoptionListingInput struct {
	UserID            string `json:"user_id"`
	GuardianProfileID string `json:"guardian_profile_id"`
	CreateAdoptionListingInputBody
}

type CreateAdoptionListingOutput struct {
	Message string                    `json:"message"`
	Listing *entities.AdoptionListing `json:"listing"`
}

type CreateAdoptionListingUseCase struct {
	petRepository             entities.PetRepository
	adoptionListingRepository entities.AdoptionListingRepository
	logger                    logging.LoggerInterface
}

func NewCreateAdoptionListingUseCase(
	petRepository entities.PetRepository,
	adoptionListingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *CreateAdoptionListingUseCase {
	return &CreateAdoptionListingUseCase{
		petRepository:             petRepository,
		adoptionListingRepository: adoptionListingRepository,
		logger:                    logger,
	}
}

func (uc *CreateAdoptionListingUseCase) Execute(ctx context.Context, input CreateAdoptionListingInput) (*CreateAdoptionListingOutput, []exceptions.ProblemDetails) {
	const from = "CreateAdoptionListingUseCase.Execute"

	if input.GuardianProfileID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do perfil de responsável ausente", errors.New("O ID do perfil de responsável é obrigatório"))
	}

	pet, err := uc.petRepository.FindByID(input.PetID)
	if err != nil {
		if err.Error() == consts.PetNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Pet não encontrado", errors.New("O pet informado não foi encontrado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar pet", err)
	}

	if pet.UserID != input.UserID {
		return nil, uc.logger.LogForbidden(ctx, from, "Pet não pertence ao usuário", errors.New("Você não tem permissão para criar um anúncio com este pet"))
	}

	listing, problems := entities.NewAdoptionListing(
		input.PetID,
		input.GuardianProfileID,
		input.Title,
		input.Description,
		input.AdoptionReason,
		input.Sex,
		input.Size,
		input.AgeGroup,
		input.CityID,
		input.StateID,
	)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Dados do anúncio inválidos", problems)
		return nil, problems
	}

	listing.Latitude = input.Latitude
	listing.Longitude = input.Longitude
	listing.Vaccinated = input.Vaccinated
	listing.Neutered = input.Neutered
	listing.Dewormed = input.Dewormed
	listing.SpecialNeeds = input.SpecialNeeds
	listing.GoodWithChildren = input.GoodWithChildren
	listing.GoodWithDogs = input.GoodWithDogs
	listing.GoodWithCats = input.GoodWithCats
	listing.RequiresHouseScreening = input.RequiresHouseScreening

	if err := uc.adoptionListingRepository.Create(listing); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar anúncio de adoção", err)
	}

	return &CreateAdoptionListingOutput{
		Message: "Anúncio de adoção criado com sucesso",
		Listing: listing,
	}, nil
}
