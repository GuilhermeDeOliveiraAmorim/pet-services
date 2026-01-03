package request

import "fmt"

// Erros de domínio
var (
	ErrRequestNotFound         = fmt.Errorf("solicitação não encontrada")
	ErrInvalidStatusTransition = fmt.Errorf("transição de status inválida")
	ErrRequestAlreadyProcessed = fmt.Errorf("solicitação já processada")
	ErrInvalidPreferredDate    = fmt.Errorf("data preferencial inválida")
	ErrInvalidServiceType      = fmt.Errorf("tipo de serviço inválido")
)
