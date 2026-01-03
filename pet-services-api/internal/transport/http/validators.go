package http

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateRequest executa validação estrutural usando tags do validator.
func ValidateRequest(req interface{}) error {
	return validate.Struct(req)
}
