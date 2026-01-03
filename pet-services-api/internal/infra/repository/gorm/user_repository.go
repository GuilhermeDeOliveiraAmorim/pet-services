package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userdom "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/models"
)

// UserRepository implementa user.Repository com GORM.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *userdom.User) error {
	mu, err := toModelUser(u)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(mu).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*userdom.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainUser(&m)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*userdom.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return toDomainUser(&m)
}

func (r *UserRepository) Update(ctx context.Context, u *userdom.User) error {
	mu, err := toModelUser(u)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(mu).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) CreatePasswordResetToken(ctx context.Context, token *userdom.PasswordResetToken) error {
	mt := toModelPasswordResetToken(token)
	return r.db.WithContext(ctx).Create(mt).Error
}

func (r *UserRepository) FindPasswordResetToken(ctx context.Context, token string) (*userdom.PasswordResetToken, error) {
	var m models.PasswordResetToken
	if err := r.db.WithContext(ctx).First(&m, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return toDomainPasswordResetToken(&m), nil
}

func (r *UserRepository) MarkPasswordResetTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("used", true).Error
}

func (r *UserRepository) CreateEmailVerificationToken(ctx context.Context, token *userdom.EmailVerificationToken) error {
	mt := toModelEmailVerificationToken(token)
	return r.db.WithContext(ctx).Create(mt).Error
}

func (r *UserRepository) FindEmailVerificationToken(ctx context.Context, token string) (*userdom.EmailVerificationToken, error) {
	var m models.EmailVerificationToken
	if err := r.db.WithContext(ctx).First(&m, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return toDomainEmailVerificationToken(&m), nil
}

func (r *UserRepository) MarkEmailVerificationTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.EmailVerificationToken{}).
		Where("id = ?", tokenID).
		Update("used", true).Error
}
