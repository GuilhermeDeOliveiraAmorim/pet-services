package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	FindByID(id string) (*RefreshToken, error)
	FindByToken(token string) (*RefreshToken, error)
	FindActiveByUserID(userID string) ([]*RefreshToken, error)
	Update(token *RefreshToken) error
	Revoke(tokenID string) error
	RevokeAllByUserID(userID string) error
	DeleteExpired() error
	IsValid(token string) (bool, error)
	CreatePasswordReset(token *PasswordResetToken) error
	RevokeAllPasswordResetByUserID(userID string) error
	FindValidPasswordResetByToken(token string) (*PasswordResetToken, error)
	RevokePasswordResetByToken(token string) error
}

type PasswordResetToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	UserAgent string
	IP        string
	RevokedAt *time.Time
}

type RefreshToken struct {
	Base
	UserID    string     `json:"user_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	UserAgent string     `json:"user_agent"`
	IpAddress string     `json:"ip_address"`
}

func NewRefreshToken(userID, token string, expiresAt time.Time, userAgent, ipAddress string) (*RefreshToken, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if token == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Token ausente",
			Detail: "O token é obrigatório",
		}))
	}

	if expiresAt.Before(time.Now()) {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Data de expiração inválida",
			Detail: "A data de expiração deve ser no futuro",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &RefreshToken{
		Base:      *NewBase(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		UserAgent: userAgent,
		IpAddress: ipAddress,
	}, nil
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

func (rt *RefreshToken) Revoke() {
	timeNow := time.Now()
	rt.RevokedAt = &timeNow
	rt.UpdatedAt = &timeNow
	rt.Deactivate()
}

func (rt *RefreshToken) IsValid() bool {
	return rt.Active && !rt.IsExpired() && !rt.IsRevoked()
}

func (rt *RefreshToken) RemainingTime() time.Duration {
	if rt.IsExpired() {
		return 0
	}
	return time.Until(rt.ExpiresAt)
}
