package entities

import "pet-services-api/internal/exceptions"

type AdoptionApplicationEvent struct {
	Base
	ApplicationID string `json:"application_id"`
	EventType     string `json:"event_type"`
	ActorUserID   string `json:"actor_user_id"`
	PayloadJSON   string `json:"payload_json"`
}

type AdoptionApplicationEventRepository interface {
	Create(event *AdoptionApplicationEvent) error
	ListByApplicationID(applicationID string) ([]*AdoptionApplicationEvent, error)
}

func NewAdoptionApplicationEvent(applicationID, eventType, actorUserID, payloadJSON string) (*AdoptionApplicationEvent, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if applicationID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da candidatura ausente",
			Detail: "O ID da candidatura é obrigatório",
		}))
	}

	if eventType == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Tipo de evento ausente",
			Detail: "O tipo de evento é obrigatório",
		}))
	} else if len(eventType) > 100 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Tipo de evento muito longo",
			Detail: "O tipo de evento deve ter no máximo 100 caracteres",
		}))
	}

	if len(payloadJSON) > 20000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Payload do evento muito longo",
			Detail: "O payload do evento deve ter no máximo 20000 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &AdoptionApplicationEvent{
		Base:          *NewBase(),
		ApplicationID: applicationID,
		EventType:     eventType,
		ActorUserID:   actorUserID,
		PayloadJSON:   payloadJSON,
	}, nil
}
