package auth

import "errors"

var (
	ErrInvalidCredentials   = errors.New("credenciais inválidas")
	ErrRefreshTokenNotFound = errors.New("refresh token não encontrado")
	ErrRefreshTokenRevoked  = errors.New("refresh token revogado")
	ErrRefreshTokenExpired  = errors.New("refresh token expirado")
)
