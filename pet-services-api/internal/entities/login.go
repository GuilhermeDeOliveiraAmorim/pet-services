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

func (l *Login) CompareAndDecrypt(hashedData string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	return err == nil
}

func (l *Login) EncryptPassword() error {
	hashedPassword, err := hashString(l.Password)
	if err != nil {
		return err
	}

	l.Password = string(hashedPassword)

	return nil
}

func (l *Login) DecryptPassword(password string) bool {
	return l.CompareAndDecrypt(l.Password, password)
}

func (l *Login) SetEmail(newEmail string) {
	l.Email = newEmail
}

func (l *Login) SetPassword(rawPassword string) {
	l.Password = rawPassword
	l.EncryptPassword()
}

func (l *Login) ChangePassword(currentPassword, newPassword string) bool {
	if !l.DecryptPassword(currentPassword) {
		return false
	}
	l.SetPassword(newPassword)
	return true
}

func (l *Login) Equals(other *Login) bool {
	return l.Email == other.Email && l.Password == other.Password
}
