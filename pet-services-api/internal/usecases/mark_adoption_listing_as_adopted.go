package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
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
	listingRepository     entities.AdoptionListingRepository
	applicationRepository entities.AdoptionApplicationRepository
	guardianRepository    entities.AdoptionGuardianProfileRepository
	userRepository        entities.UserRepository
	emailService          mail.EmailService
	logger                logging.LoggerInterface
}

func NewMarkAdoptionListingAsAdoptedUseCase(
	listingRepository entities.AdoptionListingRepository,
	applicationRepository entities.AdoptionApplicationRepository,
	guardianRepository entities.AdoptionGuardianProfileRepository,
	userRepository entities.UserRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *MarkAdoptionListingAsAdoptedUseCase {
	return &MarkAdoptionListingAsAdoptedUseCase{
		listingRepository:     listingRepository,
		applicationRepository: applicationRepository,
		guardianRepository:    guardianRepository,
		userRepository:        userRepository,
		emailService:          emailService,
		logger:                logger,
	}
}

func (u *MarkAdoptionListingAsAdoptedUseCase) Execute(ctx context.Context, input MarkAdoptionListingAsAdoptedInput) (*MarkAdoptionListingAsAdoptedOutput, []exceptions.ProblemDetails) {
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

	if listing.GuardianProfileID != input.GuardianProfileID {
		return nil, u.logger.LogForbidden(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Permissão negada", errors.New("Você não tem permissão para marcar este anúncio como adotado"))
	}

	listing.MarkAdopted()

	if err := u.listingRepository.Update(listing); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao atualizar anúncio",
			Detail: "Ocorreu um erro ao marcar o anúncio como adotado",
		})
		u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Anúncio "+input.ListingID+" marcado como adotado")

	go func() {

		guardian, err := u.guardianRepository.FindByID(listing.GuardianProfileID)
		if err == nil && guardian != nil {
			guardianUser, err := u.userRepository.FindByID(guardian.UserID)
			if err == nil && guardianUser != nil {
				if err := u.emailService.SendPetAdoptedGuardianEmail(guardianUser.Login.Email, guardianUser.Name, listing.Title); err != nil {
					u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Erro ao enviar email ao guardian", err)
				}
			}
		}

		applications, _, err := u.applicationRepository.ListByListingID(listing.ID, 1, 100)
		if err == nil && applications != nil {
			for _, app := range applications {
				applicantUser, err := u.userRepository.FindByID(app.ApplicantUserID)
				if err == nil && applicantUser != nil {
					if app.Status == "accepted" {

						if err := u.emailService.SendPetAdoptedApplicantEmail(applicantUser.Login.Email, applicantUser.Name, listing.Title); err != nil {
							u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Erro ao enviar email ao adotante", err)
						}
					} else if app.Status != "rejected" && app.Status != "withdrawn" {

						if err := u.emailService.SendPetAdoptedRejectedApplicantsEmail(applicantUser.Login.Email, listing.Title); err != nil {
							u.logger.LogError(ctx, "MarkAdoptionListingAsAdoptedUseCase", "Erro ao enviar email aos candidatos não selecionados", err)
						}
					}
				}
			}
		}
	}()

	return &MarkAdoptionListingAsAdoptedOutput{
		ID:     listing.ID,
		Status: listing.Status,
	}, nil
}
