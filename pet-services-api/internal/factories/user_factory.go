package factories

import (
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type UserFactory struct {
	RegisterUser        *usecases.RegisterUserUseCase
	RegisterAdmin       *usecases.CreateAdminUseCase
	GetProfile          *usecases.GetProfileUseCase
	ListUsers           *usecases.ListUsersUseCase
	UpdateUser          *usecases.UpdateUserUseCase
	DeleteUser          *usecases.DeleteUserUseCase
	ReactivateUser      *usecases.ReactivateUserUseCase
	DeactivateUser      *usecases.DeactivateUserUseCase
	GetUserByID         *usecases.GetUserByIDUseCase
	CheckEmailExists    *usecases.CheckEmailExistsUseCase
	CheckPhoneExists    *usecases.CheckPhoneExistsUseCase
	UpdateEmailVerified *usecases.UpdateEmailVerifiedUseCase
	ChangePassword      *usecases.ChangePasswordUseCase
}

func NewUserFactory(db *gorm.DB) *UserFactory {
	userRepo := repository_impl.NewUserRepository(db)
	tokenRepo := repository_impl.NewRefreshTokenRepository(db)

	return &UserFactory{
		RegisterUser:        usecases.NewRegisterUserUseCase(userRepo),
		RegisterAdmin:       usecases.NewCreateAdminUseCase(userRepo),
		GetProfile:          usecases.NewGetProfileUseCase(userRepo),
		ListUsers:           usecases.NewListUsersUseCase(userRepo),
		UpdateUser:          usecases.NewUpdateUserUseCase(userRepo),
		DeleteUser:          usecases.NewDeleteUserUseCase(userRepo),
		ReactivateUser:      usecases.NewReactivateUserUseCase(userRepo),
		DeactivateUser:      usecases.NewDeactivateUserUseCase(userRepo, tokenRepo),
		GetUserByID:         usecases.NewGetUserByIDUseCase(userRepo),
		CheckEmailExists:    usecases.NewCheckEmailExistsUseCase(userRepo),
		CheckPhoneExists:    usecases.NewCheckPhoneExistsUseCase(userRepo),
		UpdateEmailVerified: usecases.NewUpdateEmailVerifiedUseCase(userRepo),
		ChangePassword:      usecases.NewChangePasswordUseCase(userRepo),
	}
}
