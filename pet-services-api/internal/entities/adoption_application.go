package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type AdoptionApplicationStatusENUM struct {
	Submitted   string `json:"submitted"`
	UnderReview string `json:"under_review"`
	Interview   string `json:"interview"`
	Approved    string `json:"approved"`
	Rejected    string `json:"rejected"`
	Withdrawn   string `json:"withdrawn"`
}

var AdoptionApplicationStatuses = AdoptionApplicationStatusENUM{
	Submitted:   "submitted",
	UnderReview: "under_review",
	Interview:   "interview",
	Approved:    "approved",
	Rejected:    "rejected",
	Withdrawn:   "withdrawn",
}

type AdoptionApplication struct {
	Base
	ListingID       string     `json:"listing_id"`
	ApplicantUserID string     `json:"applicant_user_id"`
	Status          string     `json:"status"`
	Motivation      string     `json:"motivation"`
	HousingType     string     `json:"housing_type"`
	HasOtherPets    bool       `json:"has_other_pets"`
	PetExperience   string     `json:"pet_experience"`
	FamilyMembers   int        `json:"family_members"`
	AgreesHomeVisit bool       `json:"agrees_home_visit"`
	ContactPhone    string     `json:"contact_phone"`
	NotesInternal   string     `json:"notes_internal"`
	ReviewedBy      string     `json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
}

type AdoptionApplicationRepository interface {
	Create(application *AdoptionApplication) error
	FindByID(id string) (*AdoptionApplication, error)
	FindByListingIDAndApplicantUserID(listingID, applicantUserID string) (*AdoptionApplication, error)
	Update(application *AdoptionApplication) error
	ListByApplicantUserID(applicantUserID string, page, pageSize int) ([]*AdoptionApplication, int64, error)
	ListByListingID(listingID string, page, pageSize int) ([]*AdoptionApplication, int64, error)
}

func NewAdoptionApplication(listingID, applicantUserID, motivation, housingType, petExperience, contactPhone string, familyMembers int, agreesHomeVisit, hasOtherPets bool) (*AdoptionApplication, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if listingID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do anúncio ausente",
			Detail: "O ID do anúncio é obrigatório",
		}))
	}

	if applicantUserID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do candidato ausente",
			Detail: "O ID do candidato é obrigatório",
		}))
	}

	if motivation == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Motivação ausente",
			Detail: "A motivação para adoção é obrigatória",
		}))
	} else if len(motivation) > 3000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Motivação muito longa",
			Detail: "A motivação deve ter no máximo 3000 caracteres",
		}))
	}

	if len(housingType) > 100 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Tipo de moradia muito longo",
			Detail: "O tipo de moradia deve ter no máximo 100 caracteres",
		}))
	}

	if len(petExperience) > 2000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Experiência com pets muito longa",
			Detail: "A experiência com pets deve ter no máximo 2000 caracteres",
		}))
	}

	if familyMembers < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Quantidade de membros da família inválida",
			Detail: "A quantidade de membros da família não pode ser negativa",
		}))
	}

	if len(contactPhone) > 30 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Telefone de contato muito longo",
			Detail: "O telefone de contato deve ter no máximo 30 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &AdoptionApplication{
		Base:            *NewBase(),
		ListingID:       listingID,
		ApplicantUserID: applicantUserID,
		Status:          AdoptionApplicationStatuses.Submitted,
		Motivation:      motivation,
		HousingType:     housingType,
		HasOtherPets:    hasOtherPets,
		PetExperience:   petExperience,
		FamilyMembers:   familyMembers,
		AgreesHomeVisit: agreesHomeVisit,
		ContactPhone:    contactPhone,
	}, nil
}

func (a *AdoptionApplication) MoveToUnderReview(reviewedBy string) {
	timeNow := time.Now()
	a.Status = AdoptionApplicationStatuses.UnderReview
	a.ReviewedBy = reviewedBy
	a.ReviewedAt = &timeNow
	a.UpdatedAt = &timeNow
}

func (a *AdoptionApplication) MoveToInterview(reviewedBy string) {
	timeNow := time.Now()
	a.Status = AdoptionApplicationStatuses.Interview
	a.ReviewedBy = reviewedBy
	a.ReviewedAt = &timeNow
	a.UpdatedAt = &timeNow
}

func (a *AdoptionApplication) Approve(reviewedBy string) {
	timeNow := time.Now()
	a.Status = AdoptionApplicationStatuses.Approved
	a.ReviewedBy = reviewedBy
	a.ReviewedAt = &timeNow
	a.UpdatedAt = &timeNow
}

func (a *AdoptionApplication) Reject(reviewedBy, notesInternal string) {
	timeNow := time.Now()
	a.Status = AdoptionApplicationStatuses.Rejected
	a.ReviewedBy = reviewedBy
	a.ReviewedAt = &timeNow
	a.NotesInternal = notesInternal
	a.UpdatedAt = &timeNow
}

func (a *AdoptionApplication) Withdraw() {
	timeNow := time.Now()
	a.Status = AdoptionApplicationStatuses.Withdrawn
	a.UpdatedAt = &timeNow
}
