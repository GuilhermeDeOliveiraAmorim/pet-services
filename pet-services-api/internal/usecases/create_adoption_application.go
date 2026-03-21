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

type CreateAdoptionApplicationInputBody struct {
	ListingID       string `json:"listing_id" binding:"required"`
	Motivation      string `json:"motivation" binding:"required"`
	HousingType     string `json:"housing_type,omitempty"`
	HasOtherPets    bool   `json:"has_other_pets,omitempty"`
	PetExperience   string `json:"pet_experience,omitempty"`
	FamilyMembers   int    `json:"family_members,omitempty"`
	AgreesHomeVisit bool   `json:"agrees_home_visit,omitempty"`
	ContactPhone    string `json:"contact_phone,omitempty"`
}

type CreateAdoptionApplicationInput struct {
	UserID string
	CreateAdoptionApplicationInputBody
}

type CreateAdoptionApplicationOutput struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	ListingID       string `json:"listing_id"`
	ApplicantUserID string `json:"applicant_user_id"`
}

type CreateAdoptionApplicationUseCase struct {
	listingRepository  entities.AdoptionListingRepository
	guardianRepository entities.AdoptionGuardianProfileRepository
	userRepository     entities.UserRepository
	applicationRepo    entities.AdoptionApplicationRepository
	emailService       mail.EmailService
	logger             logging.LoggerInterface
}

func NewCreateAdoptionApplicationUseCase(
	listingRepository entities.AdoptionListingRepository,
	guardianRepository entities.AdoptionGuardianProfileRepository,
	userRepository entities.UserRepository,
	applicationRepo entities.AdoptionApplicationRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *CreateAdoptionApplicationUseCase {
	return &CreateAdoptionApplicationUseCase{
		listingRepository:  listingRepository,
		guardianRepository: guardianRepository,
		userRepository:     userRepository,
		applicationRepo:    applicationRepo,
		emailService:       emailService,
		logger:             logger,
	}
}

func (u *CreateAdoptionApplicationUseCase) Execute(ctx context.Context, input CreateAdoptionApplicationInput) (*CreateAdoptionApplicationOutput, []exceptions.ProblemDetails) {
	// Validar se anúncio existe e está publicado
	listing, err := u.listingRepository.FindByID(input.ListingID)
	if err != nil {
		if errors.Is(err, errors.New(consts.AdoptionListingNotFoundError)) {
			return nil, u.logger.LogNotFound(ctx, "CreateAdoptionApplicationUseCase", "Anúncio não encontrado", errors.New("O anúncio especificado não existe"))
		}
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao buscar anúncio",
			Detail: "Ocorreu um erro ao buscar o anúncio",
		})
		u.logger.LogError(ctx, "CreateAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Verificar se anúncio está publicado
	if listing.Status != entities.AdoptionListingStatuses.Published {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Anúncio não disponível",
			Detail: "Este anúncio não está disponível para candidaturas",
		})
		u.logger.LogBadRequest(ctx, "CreateAdoptionApplicationUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Verificar se usuário já se candidatou
	_, err = u.applicationRepo.FindByListingIDAndApplicantUserID(input.ListingID, input.UserID)
	if err == nil {
		// Já existe candidatura
		problem := exceptions.NewProblemDetails(exceptions.Conflict, exceptions.ErrorMessage{
			Title:  "Candidatura duplicada",
			Detail: "Você já se candidatou a este anúncio",
		})
		u.logger.LogBadRequest(ctx, "CreateAdoptionApplicationUseCase", problem.Detail, nil)
		return nil, []exceptions.ProblemDetails{problem}
	}

	// Criar candidatura
	application, problems := entities.NewAdoptionApplication(
		input.ListingID,
		input.UserID,
		input.Motivation,
		input.HousingType,
		input.PetExperience,
		input.ContactPhone,
		input.FamilyMembers,
		input.AgreesHomeVisit,
		input.HasOtherPets,
	)

	if len(problems) > 0 {
		u.logger.LogBadRequest(ctx, "CreateAdoptionApplicationUseCase", "Validação falhou", nil)
		return nil, problems
	}

	// Persistir
	if err := u.applicationRepo.Create(application); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao criar candidatura",
			Detail: "Ocorreu um erro ao criar a candidatura",
		})
		u.logger.LogError(ctx, "CreateAdoptionApplicationUseCase", problem.Detail, err)
		return nil, []exceptions.ProblemDetails{problem}
	}

	u.logger.LogInfo(ctx, "CreateAdoptionApplicationUseCase", "Candidatura criada: "+application.ID)

	// Buscar dados para envio de emails
	applicantUser, err := u.userRepository.FindByID(input.UserID)
	guardian, err2 := u.guardianRepository.FindByID(listing.GuardianProfileID)

	if err == nil && applicantUser != nil {
		// Enviar email ao candidato
		if err := u.emailService.SendAdoptionApplicationSubmittedEmail(applicantUser.Login.Email, applicantUser.Name, listing.Title); err != nil {
			u.logger.LogError(ctx, "CreateAdoptionApplicationUseCase", "Erro ao enviar email ao candidato", err)
		}
	}

	if err2 == nil && guardian != nil {
		// Buscar email do guardian
		guardianUser, err := u.userRepository.FindByID(guardian.UserID)
		if err == nil && guardianUser != nil && applicantUser != nil {
			// Enviar email ao guardian
			if err := u.emailService.SendAdoptionApplicationReceivedGuardianEmail(guardianUser.Login.Email, guardian.DisplayName, applicantUser.Name, listing.Title); err != nil {
				u.logger.LogError(ctx, "CreateAdoptionApplicationUseCase", "Erro ao enviar email ao guardian", err)
			}
		}
	}

	return &CreateAdoptionApplicationOutput{
		ID:              application.ID,
		Status:          application.Status,
		ListingID:       application.ListingID,
		ApplicantUserID: application.ApplicantUserID,
	}, nil
}
