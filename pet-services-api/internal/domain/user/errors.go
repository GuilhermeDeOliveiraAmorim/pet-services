package user

import "fmt"

// Erros de domínio
var (
	ErrUserNotFound      = fmt.Errorf("usuário não encontrado")
	ErrUserAlreadyExists = fmt.Errorf("usuário já existe")
	ErrInvalidEmail      = fmt.Errorf("email inválido")
	ErrInvalidPhone      = fmt.Errorf("telefone inválido")
	ErrInvalidUserType   = fmt.Errorf("tipo de usuário inválido")
	ErrInvalidPassword   = fmt.Errorf("senha inválida")
)
