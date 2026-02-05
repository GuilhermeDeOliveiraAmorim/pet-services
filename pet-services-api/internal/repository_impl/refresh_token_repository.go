package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) entities.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *entities.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByID(id string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken
	err := r.db.First(&token, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) FindByToken(tokenStr string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken
	err := r.db.First(&token, "token = ?", tokenStr).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) FindActiveByUserID(userID string) ([]*entities.RefreshToken, error) {
	var tokens []*entities.RefreshToken
	err := r.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *refreshTokenRepository) Update(token *entities.RefreshToken) error {
	return r.db.Save(token).Error
}

func (r *refreshTokenRepository) Revoke(tokenID string) error {
	return r.db.Model(&entities.RefreshToken{}).Where("id = ?", tokenID).Update("revoked_at", time.Now()).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(userID string) error {
	return r.db.Model(&entities.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", time.Now()).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&entities.RefreshToken{}).Error
}

func (r *refreshTokenRepository) IsValid(token string) (bool, error) {
	t, err := r.FindByToken(token)
	if err != nil {
		return false, err
	}
	return t.IsValid(), nil
}

func (r *refreshTokenRepository) CreatePasswordReset(token *entities.PasswordResetToken) error {
	model := &models.PasswordResetToken{}
	model.FromEntity(token)
	return r.db.Create(model).Error
}

func (r *refreshTokenRepository) RevokeAllPasswordResetByUserID(userID string) error {
	return r.db.Model(&models.PasswordResetToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", time.Now()).Error
}

func (r *refreshTokenRepository) FindValidPasswordResetByToken(token string) (*entities.PasswordResetToken, error) {
	var model models.PasswordResetToken
	err := r.db.Where("token = ? AND revoked_at IS NULL AND expires_at > ?", token, time.Now()).First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *refreshTokenRepository) RevokePasswordResetByToken(token string) error {
	return r.db.Model(&models.PasswordResetToken{}).Where("token = ? AND revoked_at IS NULL", token).Update("revoked_at", time.Now()).Error
}
