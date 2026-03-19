package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
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

type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
	Issuer          string
}

type JWTService struct {
	cfg    JWTConfig
	parser *jwt.Parser
}

func NewJWTServiceFromEnv() (*JWTService, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET não configurado")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = accessSecret
	}

	accessDur := mustParseDuration(os.Getenv("JWT_ACCESS_DURATION"), 15*time.Minute)
	refreshDur := mustParseDuration(os.Getenv("JWT_REFRESH_DURATION"), 168*time.Hour)

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "pet-services"
	}

	return &JWTService{
		cfg: JWTConfig{
			AccessSecret:    accessSecret,
			RefreshSecret:   refreshSecret,
			AccessDuration:  accessDur,
			RefreshDuration: refreshDur,
			Issuer:          issuer,
		},
		parser: jwt.NewParser(),
	}, nil
}

func (s *JWTService) GenerateAccessToken(userID, email, userType string) (string, error) {
	return s.generateToken(userID, email, userType, AccessToken, s.cfg.AccessDuration)
}

func (s *JWTService) GenerateRefreshToken(userID, email, userType string) (string, error) {
	return s.generateToken(userID, email, userType, RefreshToken, s.cfg.RefreshDuration)
}

func (s *JWTService) GenerateTokenPair(userID, email, userType string) (*TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(userID, email, userType)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(userID, email, userType)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.AccessDuration.Seconds()),
	}, nil
}

func (s *JWTService) generateToken(userID, email, userType string, tokenType TokenType, duration time.Duration) (string, error) {
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
			Issuer:    s.cfg.Issuer,
			ID:        ulid.Make().String(),
		},
	}

	secret := s.getSecret(tokenType)
	if secret == "" {
		return "", errors.New("secret JWT não configurado")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *JWTService) ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	secret := s.getSecret(tokenType)
	if secret == "" {
		return nil, errors.New("secret JWT não configurado")
	}

	token, err := s.parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.ValidateToken(tokenString, AccessToken)
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.ValidateToken(tokenString, RefreshToken)
}

func (s *JWTService) ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := s.parser.ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JWTService) getSecret(tokenType TokenType) string {
	if tokenType == RefreshToken {
		return s.cfg.RefreshSecret
	}
	return s.cfg.AccessSecret
}

func mustParseDuration(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}

	dur, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return dur
}

func (s *JWTService) IsTokenExpired(tokenString string) bool {
	claims, err := s.ExtractClaims(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Time.Before(time.Now())
}
