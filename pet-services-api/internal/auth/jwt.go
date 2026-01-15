package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("token inválido")
	ErrExpiredToken     = errors.New("token expirado")
	ErrInvalidSignature = errors.New("assinatura inválida")
	ErrMissingToken     = errors.New("token ausente")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID   string    `json:"user_id"`
	Email    string    `json:"email"`
	UserType string    `json:"user_type"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func GenerateAccessToken(userID, email, userType string) (string, error) {
	durationStr := os.Getenv("JWT_ACCESS_DURATION")
	if durationStr == "" {
		durationStr = "15m"
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration = 15 * time.Minute
	}

	return generateToken(userID, email, userType, AccessToken, duration)
}

func GenerateRefreshToken(userID, email, userType string) (string, error) {
	durationStr := os.Getenv("JWT_REFRESH_DURATION")
	if durationStr == "" {
		durationStr = "168h"
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration = 168 * time.Hour
	}

	return generateToken(userID, email, userType, RefreshToken, duration)
}

func GenerateTokenPair(userID, email, userType string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, email, userType)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(userID, email, userType)
	if err != nil {
		return nil, err
	}

	durationStr := os.Getenv("JWT_ACCESS_DURATION")
	if durationStr == "" {
		durationStr = "15m"
	}
	duration, _ := time.ParseDuration(durationStr)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(duration.Seconds()),
	}, nil
}

func generateToken(userID, email, userType string, tokenType TokenType, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := Claims{
		UserID:   userID,
		Email:    email,
		UserType: userType,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "pet-services",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := getSecret(tokenType)
	if secret == "" {
		return "", errors.New("secret JWT não configurado")
	}

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	secret := getSecret(tokenType)
	if secret == "" {
		return nil, errors.New("secret JWT não configurado")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.Type != tokenType {
		return nil, errors.New("tipo de token incorreto")
	}

	return claims, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	return ValidateToken(tokenString, AccessToken)
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return ValidateToken(tokenString, RefreshToken)
}

func ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func getSecret(tokenType TokenType) string {
	if tokenType == RefreshToken {
		secret := os.Getenv("JWT_REFRESH_SECRET")
		if secret == "" {
			return os.Getenv("JWT_ACCESS_SECRET")
		}
		return secret
	}
	return os.Getenv("JWT_ACCESS_SECRET")
}

func IsTokenExpired(tokenString string) bool {
	claims, err := ExtractClaims(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Time.Before(time.Now())
}
