package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type AdoptionListingStatusENUM struct {
	Draft     string `json:"draft"`
	Published string `json:"published"`
	Paused    string `json:"paused"`
	InProcess string `json:"in_process"`
	Adopted   string `json:"adopted"`
	Archived  string `json:"archived"`
}

var AdoptionListingStatuses = AdoptionListingStatusENUM{
	Draft:     "draft",
	Published: "published",
	Paused:    "paused",
	InProcess: "in_process",
	Adopted:   "adopted",
	Archived:  "archived",
}

type AdoptionPetSexENUM struct {
	Male   string `json:"male"`
	Female string `json:"female"`
}

var AdoptionPetSexes = AdoptionPetSexENUM{
	Male:   "male",
	Female: "female",
}

type AdoptionPetSizeENUM struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

var AdoptionPetSizes = AdoptionPetSizeENUM{
	Small:  "small",
	Medium: "medium",
	Large:  "large",
}

type AdoptionPetAgeGroupENUM struct {
	Puppy  string `json:"puppy"`
	Adult  string `json:"adult"`
	Senior string `json:"senior"`
}

var AdoptionPetAgeGroups = AdoptionPetAgeGroupENUM{
	Puppy:  "puppy",
	Adult:  "adult",
	Senior: "senior",
}

type AdoptionListing struct {
	Base
	PetID                  string                  `json:"pet_id"`
	Pet                    Pet                     `json:"pet"`
	GuardianProfileID      string                  `json:"guardian_profile_id"`
	GuardianProfile        AdoptionGuardianProfile `json:"guardian_profile"`
	Title                  string                  `json:"title"`
	Description            string                  `json:"description"`
	AdoptionReason         string                  `json:"adoption_reason"`
	Status                 string                  `json:"status"`
	Sex                    string                  `json:"sex"`
	Size                   string                  `json:"size"`
	AgeGroup               string                  `json:"age_group"`
	Vaccinated             bool                    `json:"vaccinated"`
	Neutered               bool                    `json:"neutered"`
	Dewormed               bool                    `json:"dewormed"`
	SpecialNeeds           bool                    `json:"special_needs"`
	GoodWithChildren       bool                    `json:"good_with_children"`
	GoodWithDogs           bool                    `json:"good_with_dogs"`
	GoodWithCats           bool                    `json:"good_with_cats"`
	RequiresHouseScreening bool                    `json:"requires_house_screening"`
	CityID                 string                  `json:"city_id"`
	StateID                string                  `json:"state_id"`
	Latitude               float64                 `json:"latitude"`
	Longitude              float64                 `json:"longitude"`
	PublishedAt            *time.Time              `json:"published_at"`
	AdoptedAt              *time.Time              `json:"adopted_at"`
}

type AdoptionListingRepository interface {
	Create(listing *AdoptionListing) error
	FindByID(id string) (*AdoptionListing, error)
	Update(listing *AdoptionListing) error
	ListPublic(filters AdoptionListingFilters, page, pageSize int) ([]*AdoptionListing, int64, error)
	ListByGuardianProfileID(guardianProfileID string, page, pageSize int) ([]*AdoptionListing, int64, error)
}

type AdoptionListingFilters struct {
	SpeciesID string `json:"species_id"`
	Sex       string `json:"sex"`
	Size      string `json:"size"`
	AgeGroup  string `json:"age_group"`
	CityID    string `json:"city_id"`
	StateID   string `json:"state_id"`
	Status    string `json:"status"`
}

func NewAdoptionListing(petID, guardianProfileID, title, description, adoptionReason, sex, size, ageGroup, cityID, stateID string) (*AdoptionListing, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if petID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		}))
	}

	if guardianProfileID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do responsável ausente",
			Detail: "O ID do responsável é obrigatório",
		}))
	}

	if title == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Título ausente",
			Detail: "O título do anúncio é obrigatório",
		}))
	} else if len(title) > 140 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Título muito longo",
			Detail: "O título do anúncio deve ter no máximo 140 caracteres",
		}))
	}

	if description == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição ausente",
			Detail: "A descrição do anúncio é obrigatória",
		}))
	} else if len(description) > 5000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição muito longa",
			Detail: "A descrição do anúncio deve ter no máximo 5000 caracteres",
		}))
	}

	if len(adoptionReason) > 1000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Motivo de adoção muito longo",
			Detail: "O motivo de adoção deve ter no máximo 1000 caracteres",
		}))
	}

	if sex != "" && sex != AdoptionPetSexes.Male && sex != AdoptionPetSexes.Female {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Sexo inválido",
			Detail: "O sexo deve ser 'male' ou 'female'",
		}))
	}

	if size != "" && size != AdoptionPetSizes.Small && size != AdoptionPetSizes.Medium && size != AdoptionPetSizes.Large {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Porte inválido",
			Detail: "O porte deve ser 'small', 'medium' ou 'large'",
		}))
	}

	if ageGroup != "" && ageGroup != AdoptionPetAgeGroups.Puppy && ageGroup != AdoptionPetAgeGroups.Adult && ageGroup != AdoptionPetAgeGroups.Senior {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Faixa etária inválida",
			Detail: "A faixa etária deve ser 'puppy', 'adult' ou 'senior'",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &AdoptionListing{
		Base:              *NewBase(),
		PetID:             petID,
		GuardianProfileID: guardianProfileID,
		Title:             title,
		Description:       description,
		AdoptionReason:    adoptionReason,
		Status:            AdoptionListingStatuses.Draft,
		Sex:               sex,
		Size:              size,
		AgeGroup:          ageGroup,
		CityID:            cityID,
		StateID:           stateID,
	}, nil
}

func (l *AdoptionListing) Publish() {
	timeNow := time.Now()
	l.Status = AdoptionListingStatuses.Published
	l.PublishedAt = &timeNow
	l.UpdatedAt = &timeNow
}

func (l *AdoptionListing) Pause() {
	timeNow := time.Now()
	l.Status = AdoptionListingStatuses.Paused
	l.UpdatedAt = &timeNow
}

func (l *AdoptionListing) MarkInProcess() {
	timeNow := time.Now()
	l.Status = AdoptionListingStatuses.InProcess
	l.UpdatedAt = &timeNow
}

func (l *AdoptionListing) MarkAdopted() {
	timeNow := time.Now()
	l.Status = AdoptionListingStatuses.Adopted
	l.AdoptedAt = &timeNow
	l.UpdatedAt = &timeNow
}

func (l *AdoptionListing) Archive() {
	timeNow := time.Now()
	l.Status = AdoptionListingStatuses.Archived
	l.UpdatedAt = &timeNow
}
