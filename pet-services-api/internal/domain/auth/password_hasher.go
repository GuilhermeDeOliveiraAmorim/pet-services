package auth

// PasswordHasher abstrai hashing e comparação de senha.
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hash, password string) error
}
