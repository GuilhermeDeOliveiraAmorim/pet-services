package review

import "fmt"

// Erros de domínio
var (
	ErrReviewNotFound      = fmt.Errorf("avaliação não encontrada")
	ErrReviewAlreadyExists = fmt.Errorf("avaliação já existe para esta solicitação")
	ErrInvalidRating       = fmt.Errorf("nota inválida")
	ErrRequestNotCompleted = fmt.Errorf("solicitação não foi concluída")
)
