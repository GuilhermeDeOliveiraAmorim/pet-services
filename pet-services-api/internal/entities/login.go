package entities

import (
	"pet-services-api/internal/exceptions"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(email, password string) (*Login, []exceptions.ProblemDetails) {
	validationErrors := ValidateLogin(email, password)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Login{
		Email:    email,
		Password: password,
	}, nil
}

func ValidateLogin(email, password string) []exceptions.ProblemDetails {
	var validationErrors []exceptions.ProblemDetails

	if !IsValidEmail(email) {
		validationErrors = append(validationErrors, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "E-mail inválido",
			Detail: "O formato do e-mail é inválido",
		}))
	}

	if !IsValidPassword(password) {
		validationErrors = append(validationErrors, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Senha inválida",
			Detail: "A senha deve ter pelo menos 6 caracteres, conter pelo menos uma letra maiúscula, uma letra minúscula, um dígito e um caractere especial",
		}))
	}

	return validationErrors
}

func IsValidEmail(email string) bool {
	emailPattern := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, _ := regexp.MatchString(emailPattern, email)
	return match
}

func IsValidPassword(password string) bool {
	return hasMinimumLength(password, 6) &&
		hasUpperCaseLetter(password) &&
		hasLowerCaseLetter(password) &&
		hasDigit(password) &&
		hasSpecialCharacter(password)
}

func hasMinimumLength(password string, length int) bool {
	return len(password) >= length
}

func hasUpperCaseLetter(password string) bool {
	return strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func hasLowerCaseLetter(password string) bool {
	return strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
}

func hasDigit(password string) bool {
	return strings.ContainsAny(password, "0123456789")
}

func hasSpecialCharacter(password string) bool {
	specialCharacters := "@#$%&*"
	return strings.ContainsAny(password, specialCharacters)
}

func hashString(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (lo *Login) CompareAndDecrypt(hashedData string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	return err == nil
}

func (lo *Login) EncryptPassword() error {
	hashedPassword, err := hashString(lo.Password)
	if err != nil {
		return err
	}

	lo.Password = string(hashedPassword)

	return nil
}

func (lo *Login) DecryptPassword(password string) bool {
	return lo.CompareAndDecrypt(lo.Password, password)
}

func (lo *Login) SetEmail(newEmail string) {
	lo.Email = newEmail
}

func (lo *Login) SetPassword(rawPassword string) {
	lo.Password = rawPassword
	lo.EncryptPassword()
}

func (lo *Login) ChangePassword(currentPassword, newPassword string) bool {
	if !lo.DecryptPassword(currentPassword) {
		return false
	}
	lo.SetPassword(newPassword)
	return true
}

func (lo *Login) Equals(other *Login) bool {
	return lo.Email == other.Email && lo.Password == other.Password
}
