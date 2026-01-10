package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"pet-services-api/internal/domain/auth"
	"pet-services-api/internal/domain/user"
)

// Config holds JWT configuration.
type Config struct {
	AccessSecret    string
	RefreshSecret   string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

// TokenService implements auth.TokenService using JWT (HS256).
type TokenService struct {
	cfg Config
}

// NewTokenService creates a new JWT token service.
func NewTokenService(cfg Config) *TokenService {
	return &TokenService{cfg: cfg}
}

// accessClaims represents JWT claims for access tokens.
type accessClaims struct {
	UserID   string        `json:"user_id"`
	UserType user.UserType `json:"user_type"`
	jwtlib.RegisteredClaims
}

// refreshClaims represents JWT claims for refresh tokens.
type refreshClaims struct {
	TokenID  uuid.UUID     `json:"token_id"`
	UserID   uuid.UUID     `json:"user_id"`
	UserType user.UserType `json:"user_type"`
	jwtlib.RegisteredClaims
}

// GenerateTokens generates access and refresh tokens.
func (s *TokenService) GenerateTokens(userID uuid.UUID, userType user.UserType) (auth.TokenPair, error) {
	now := time.Now()
	accessExpiresAt := now.Add(s.cfg.AccessDuration)
	refreshExpiresAt := now.Add(s.cfg.RefreshDuration)
	refreshID := uuid.New()

	// Generate access token
	accessToken, err := s.generateAccessToken(userID, userType, accessExpiresAt)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(refreshID, userID, userType, refreshExpiresAt)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return auth.TokenPair{
		AccessToken:      accessToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
		RefreshID:        refreshID,
	}, nil
}

// ParseRefreshToken parses and validates a refresh token, returning its claims.
func (s *TokenService) ParseRefreshToken(token string) (auth.RefreshClaims, error) {
	claims := &refreshClaims{}
	parsed, err := jwtlib.ParseWithClaims(token, claims, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.RefreshSecret), nil
	})

	if err != nil {
		return auth.RefreshClaims{}, fmt.Errorf("parse refresh token: %w", err)
	}

	if !parsed.Valid {
		return auth.RefreshClaims{}, errors.New("invalid refresh token")
	}

	return auth.RefreshClaims{
		TokenID:   claims.TokenID,
		UserID:    claims.UserID,
		UserType:  claims.UserType,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// ParseAccessToken parses and validates an access token.
// Not used in auth middleware but provided for completeness.
func (s *TokenService) ParseAccessToken(token string) (auth.AccessClaims, error) {
	claims := &accessClaims{}
	parsed, err := jwtlib.ParseWithClaims(token, claims, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			fmt.Printf("[JWT] ParseAccessToken: método de assinatura inesperado: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.AccessSecret), nil
	})

	if err != nil {
		fmt.Printf("[JWT] ParseAccessToken: erro ao fazer parse do token: %v\n", err)
		return auth.AccessClaims{}, fmt.Errorf("parse access token: %w", err)
	}

	if !parsed.Valid {
		fmt.Printf("[JWT] ParseAccessToken: token inválido\n")
		return auth.AccessClaims{}, errors.New("invalid access token")
	}

	fmt.Printf("[JWT] ParseAccessToken: claims extraídas: user_id=%v, user_type=%v, expires_at=%v\n", claims.UserID, claims.UserType, claims.ExpiresAt.Time)

	// Converte user_id string para uuid.UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		fmt.Printf("[JWT] ParseAccessToken: user_id inválido: %v\n", err)
		return auth.AccessClaims{}, fmt.Errorf("user_id inválido: %w", err)
	}

	return auth.AccessClaims{
		UserID:    userID,
		UserType:  claims.UserType,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *TokenService) generateAccessToken(userID uuid.UUID, userType user.UserType, expiresAt time.Time) (string, error) {
	claims := &accessClaims{
		UserID:   userID.String(),
		UserType: userType,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
			Issuer:    "pet-services-api",
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) generateRefreshToken(tokenID, userID uuid.UUID, userType user.UserType, expiresAt time.Time) (string, error) {
	claims := &refreshClaims{
		TokenID:  tokenID,
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
			Issuer:    "pet-services-api",
			ID:        tokenID.String(),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
