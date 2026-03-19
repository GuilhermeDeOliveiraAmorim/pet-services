package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"
	"time"

	"gorm.io/gorm"
)

type UserFactory struct {
	RegisterUser     *usecases.RegisterUserUseCase
	RegisterAdmin    *usecases.CreateAdminUseCase
	GetProfile       *usecases.GetProfileUseCase
	ListUsers        *usecases.ListUsersUseCase
	UpdateUser       *usecases.UpdateUserUseCase
	DeleteUser       *usecases.DeleteUserUseCase
	ReactivateUser   *usecases.ReactivateUserUseCase
	DeactivateUser   *usecases.DeactivateUserUseCase
	GetUserByID      *usecases.GetUserByIDUseCase
	CheckEmailExists *usecases.CheckEmailExistsUseCase
	CheckPhoneExists *usecases.CheckPhoneExistsUseCase
	ChangePassword   *usecases.ChangePasswordUseCase
	AddUserPhoto     *usecases.AddUserPhotoUseCase
}

func NewUserFactory(db *gorm.DB, storageService storage.ObjectStorage, mailService mail.EmailService, ttl time.Duration, logger logging.LoggerInterface) *UserFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	tokenRepo := repository_impl.NewRefreshTokenRepository(db)
	photoRepo := repository_impl.NewPhotoRepository(db)

	return &UserFactory{
		RegisterUser:     usecases.NewRegisterUserUseCase(userRepo, tokenRepo, mailService, ttl, logger),
		RegisterAdmin:    usecases.NewCreateAdminUseCase(userRepo, logger),
		GetProfile:       usecases.NewGetProfileUseCase(userRepo, providerRepo, storageService, logger),
		ListUsers:        usecases.NewListUsersUseCase(userRepo, storageService, logger),
		UpdateUser:       usecases.NewUpdateUserUseCase(userRepo, storageService, logger),
		DeleteUser:       usecases.NewDeleteUserUseCase(userRepo, logger),
		ReactivateUser:   usecases.NewReactivateUserUseCase(userRepo, logger),
		DeactivateUser:   usecases.NewDeactivateUserUseCase(userRepo, tokenRepo, logger),
		GetUserByID:      usecases.NewGetUserByIDUseCase(userRepo, storageService, logger),
		CheckEmailExists: usecases.NewCheckEmailExistsUseCase(userRepo, logger),
		CheckPhoneExists: usecases.NewCheckPhoneExistsUseCase(userRepo, logger),
		ChangePassword:   usecases.NewChangePasswordUseCase(userRepo, mailService, logger),
		AddUserPhoto:     usecases.NewAddUserPhotoUseCase(userRepo, photoRepo, storageService, logger),
	}
}
