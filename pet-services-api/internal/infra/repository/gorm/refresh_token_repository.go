package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authdom "github.com/guilherme/pet-services-api/internal/domain/auth"
	"github.com/guilherme/pet-services-api/internal/models"
)

// RefreshTokenRepository implementa RefreshTokenRepository com GORM.
type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *authdom.RefreshToken) error {
	mt := toModelRefreshToken(token)
	return r.db.WithContext(ctx).Create(mt).Error
}

func (r *RefreshTokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*authdom.RefreshToken, error) {
	var m models.RefreshToken
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainRefreshToken(&m), nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}
