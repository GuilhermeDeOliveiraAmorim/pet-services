package factories

import (
	"pet-services-api/internal/logging"
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

func NewUserFactory(db *gorm.DB, logger logging.LoggerInterface) *UserFactory {
	userRepo := repository_impl.NewUserRepository(db)
	tokenRepo := repository_impl.NewRefreshTokenRepository(db)

	return &UserFactory{
		RegisterUser:        usecases.NewRegisterUserUseCase(userRepo, logger),
		RegisterAdmin:       usecases.NewCreateAdminUseCase(userRepo, logger),
		GetProfile:          usecases.NewGetProfileUseCase(userRepo, logger),
		ListUsers:           usecases.NewListUsersUseCase(userRepo, logger),
		UpdateUser:          usecases.NewUpdateUserUseCase(userRepo, logger),
		DeleteUser:          usecases.NewDeleteUserUseCase(userRepo, logger),
		ReactivateUser:      usecases.NewReactivateUserUseCase(userRepo, logger),
		DeactivateUser:      usecases.NewDeactivateUserUseCase(userRepo, tokenRepo, logger),
		GetUserByID:         usecases.NewGetUserByIDUseCase(userRepo, logger),
		CheckEmailExists:    usecases.NewCheckEmailExistsUseCase(userRepo, logger),
		CheckPhoneExists:    usecases.NewCheckPhoneExistsUseCase(userRepo, logger),
		UpdateEmailVerified: usecases.NewUpdateEmailVerifiedUseCase(userRepo, logger),
		ChangePassword:      usecases.NewChangePasswordUseCase(userRepo, logger),
	}
}
