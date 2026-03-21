package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type MarkAdoptionListingAsAdoptedInput struct {
	ListingID         string
	GuardianProfileID string
}

type MarkAdoptionListingAsAdoptedOutput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type MarkAdoptionListingAsAdoptedUseCase struct {
	listingRepository entities.AdoptionListingRepository
	logger            logging.LoggerInterface
}

func NewMarkAdoptionListingAsAdoptedUseCase(
	listingRepository entities.AdoptionListingRepository,
	logger logging.LoggerInterface,
) *MarkAdoptionListingAsAdoptedUseCase {
	return &MarkAdoptionListingAsAdoptedUseCase{
		listingRepository: listingRepository,
		logger:            logger,
	}
}

func (u *MarkAdoptionListingAsAdoptedUseCase) Execute(ctx context.Context, input MarkAdoptionListingAsAdoptedInput) (*MarkAdoptionListingAsAdoptedOutput, []exceptions.ProblemDetails) {
	// Buscar anúncio
	listing, err := u.listingRepository.FindByID(input.ListingID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionListingNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Anúncio não encontrado", errors.New("O anúncio especificado não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar anúncio",
			Detail: "Ocorreu um erro ao buscar o anúncio",
		})
		u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Verificar se pertence ao perfil de guardian
	if listing.GuardianProfileID != input.GuardianProfileID {
		return nil, u.logger.LogForbidden(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Permissão negada", errors.New("Você não tem permissão para marcar este anúncio como adotado"))
	}

	// Marcar como adotado
	listing.MarkAdopted()

	// Persistir
	if err := u.listingRepository.Update(listing); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar anúncio",
			Detail: "Ocorreu um erro ao marcar o anúncio como adotado",
		})
		u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Anúncio "+input.ListingID+" marcado como adotado")

	return &MarkAdoptionListingAsAdoptedOutput{
		ID:     listing.ID,
		Status: listing.Status,
	}, nil
}
