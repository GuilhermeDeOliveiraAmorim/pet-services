package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type AdoptionGuardianTypeENUM struct {
	NGO         string `json:"ngo"`
	Independent string `json:"independent"`
	Owner       string `json:"owner"`
}

var AdoptionGuardianTypes = AdoptionGuardianTypeENUM{
	NGO:         "ngo",
	Independent: "independent",
	Owner:       "owner",
}

type AdoptionGuardianApprovalStatusENUM struct {
	Pending  string `json:"pending"`
	Approved string `json:"approved"`
	Rejected string `json:"rejected"`
}

var AdoptionGuardianApprovalStatuses = AdoptionGuardianApprovalStatusENUM{
	Pending:  "pending",
	Approved: "approved",
	Rejected: "rejected",
}

type AdoptionGuardianProfile struct {
	Base
	UserID         string     `json:"user_id"`
	DisplayName    string     `json:"display_name"`
	GuardianType   string     `json:"guardian_type"`
	Document       string     `json:"document"`
	Phone          string     `json:"phone"`
	Whatsapp       string     `json:"whatsapp"`
	About          string     `json:"about"`
	CityID         string     `json:"city_id"`
	StateID        string     `json:"state_id"`
	ApprovalStatus string     `json:"approval_status"`
	ApprovedBy     string     `json:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at"`
}

type AdoptionGuardianProfileRepository interface {
	Create(profile *AdoptionGuardianProfile) error
	FindByID(id string) (*AdoptionGuardianProfile, error)
	FindByUserID(userID string) (*AdoptionGuardianProfile, error)
	Update(profile *AdoptionGuardianProfile) error
	ListByApprovalStatus(status string, page, pageSize int) ([]*AdoptionGuardianProfile, int64, error)
}

func NewAdoptionGuardianProfile(userID, displayName, guardianType, document, phone, whatsapp, about, cityID, stateID string) (*AdoptionGuardianProfile, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if displayName == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome de exibição ausente",
			Detail: "O nome de exibição é obrigatório",
		}))
	} else if len(displayName) > 120 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome de exibição muito longo",
			Detail: "O nome de exibição deve ter no máximo 120 caracteres",
		}))
	}

	if guardianType != AdoptionGuardianTypes.NGO && guardianType != AdoptionGuardianTypes.Independent && guardianType != AdoptionGuardianTypes.Owner {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Tipo de responsável inválido",
			Detail: "O tipo de responsável deve ser 'ngo', 'independent' ou 'owner'",
		}))
	}

	if len(document) > 30 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Documento muito longo",
			Detail: "O documento deve ter no máximo 30 caracteres",
		}))
	}

	if len(phone) > 30 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Telefone muito longo",
			Detail: "O telefone deve ter no máximo 30 caracteres",
		}))
	}

	if len(whatsapp) > 30 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "WhatsApp muito longo",
			Detail: "O WhatsApp deve ter no máximo 30 caracteres",
		}))
	}

	if len(about) > 2000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição muito longa",
			Detail: "A descrição deve ter no máximo 2000 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &AdoptionGuardianProfile{
		Base:           *NewBase(),
		UserID:         userID,
		DisplayName:    displayName,
		GuardianType:   guardianType,
		Document:       document,
		Phone:          phone,
		Whatsapp:       whatsapp,
		About:          about,
		CityID:         cityID,
		StateID:        stateID,
		ApprovalStatus: AdoptionGuardianApprovalStatuses.Pending,
	}, nil
}

func (p *AdoptionGuardianProfile) Approve(approvedBy string) {
	timeNow := time.Now()
	p.ApprovalStatus = AdoptionGuardianApprovalStatuses.Approved
	p.ApprovedBy = approvedBy
	p.ApprovedAt = &timeNow
	p.UpdatedAt = &timeNow
}

func (p *AdoptionGuardianProfile) Reject(reviewedBy string) {
	timeNow := time.Now()
	p.ApprovalStatus = AdoptionGuardianApprovalStatuses.Rejected
	p.ApprovedBy = reviewedBy
	p.ApprovedAt = &timeNow
	p.UpdatedAt = &timeNow
}
