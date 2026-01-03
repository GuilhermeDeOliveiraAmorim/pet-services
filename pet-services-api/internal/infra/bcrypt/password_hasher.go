package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher implements auth.PasswordHasher using bcrypt.
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new bcrypt password hasher.
// cost should be bcrypt.DefaultCost (typically 10) or higher for security.
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash hashes a password using bcrypt.
func (h *PasswordHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compare compares a bcrypt hash with a plaintext password.
func (h *PasswordHasher) Compare(hash, password string) error {
	if hash == "" || password == "" {
		return errors.New("hash and password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
