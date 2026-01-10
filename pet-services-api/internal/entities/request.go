package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type RequestStatusENUM struct {
	Pending  string `json:"pending"`
	Accepted string `json:"accepted"`
	Rejected string `json:"rejected"`
}

var RequestStatuses = RequestStatusENUM{
	Pending:  "pending",
	Accepted: "accepted",
	Rejected: "rejected",
}

type Request struct {
	Base
	UserID       string `json:"user_id"`
	ProviderID   string `json:"provider_id"`
	ServiceID    string `json:"service_id"`
	Pet          Pet    `json:"pet"`
	Notes        string `json:"notes"`
	Status       string `json:"status"`
	RejectReason string `json:"reject_reason"`
}

func NewRequest(userID, providerID, serviceID string, pet Pet, notes string) (*Request, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if providerID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do provedor ausente",
			Detail: "O ID do provedor é obrigatório",
		}))
	}

	if serviceID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do serviço ausente",
			Detail: "O ID do serviço é obrigatório",
		}))
	}

	if len(notes) > 500 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Observações muito longas",
			Detail: "As observações devem ter no máximo 500 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Request{
		Base:         *NewBase(),
		UserID:       userID,
		ProviderID:   providerID,
		ServiceID:    serviceID,
		Pet:          pet,
		Notes:        notes,
		Status:       RequestStatuses.Pending,
		RejectReason: "",
	}, nil
}

func (r *Request) Accept() {
	timeNow := time.Now()
	r.Status = RequestStatuses.Accepted
	r.RejectReason = ""
	r.UpdatedAt = &timeNow
}

func (r *Request) Reject(reason string) {
	timeNow := time.Now()
	r.Status = RequestStatuses.Rejected
	r.RejectReason = reason
	r.UpdatedAt = &timeNow
}
