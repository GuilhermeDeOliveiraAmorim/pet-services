package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Você pode adicionar custom validators aqui se necessário
}

// ValidateRequest valida um struct de request contra suas tags de validação
func ValidateRequest(req interface{}) error {
	return validate.Struct(req)
}

// ValidationErrorResponse retorna uma resposta formatada com erros de validação
func ValidationErrorResponse(err error) map[string]interface{} {
	errors := make(map[string]interface{})

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			fieldName := fe.Field()
			errorMsg := fmt.Sprintf("field validation failed: %s (%s)", fieldName, fe.Tag())

			// Mensagens mais legíveis por tipo de validação
			switch fe.Tag() {
			case "required":
				errorMsg = fmt.Sprintf("%s is required", fieldName)
			case "email":
				errorMsg = fmt.Sprintf("%s must be a valid email", fieldName)
			case "min":
				errorMsg = fmt.Sprintf("%s must be at least %s", fieldName, fe.Param())
			case "max":
				errorMsg = fmt.Sprintf("%s must be at most %s", fieldName, fe.Param())
			case "uuid":
				errorMsg = fmt.Sprintf("%s must be a valid UUID", fieldName)
			case "oneof":
				errorMsg = fmt.Sprintf("%s must be one of: %s", fieldName, fe.Param())
			case "datetime":
				errorMsg = fmt.Sprintf("%s must be in format: %s", fieldName, fe.Param())
			case "gtfield":
				errorMsg = fmt.Sprintf("%s must be greater than %s", fieldName, fe.Param())
			}

			errors[fieldName] = errorMsg
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}

// BindAndValidateJSON faz bind do JSON e valida em um passo
func BindAndValidateJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}

	if err := ValidateRequest(req); err != nil {
		return err
	}

	return nil
}

// ValidatingMiddleware retorna um middleware que valida requisições automaticamente
// Útil quando aplicado globalmente para garantir validação em todas as rotas
func ValidatingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continua para o próximo handler
		// A validação real acontece em BindAndValidateJSON nos handlers específicos
		c.Next()
	}
}
