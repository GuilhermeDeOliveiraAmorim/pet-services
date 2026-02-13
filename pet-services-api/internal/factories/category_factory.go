package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type CategoryFactory struct {
	CreateCategory *usecases.CreateCategoryUseCase
	ListCategories *usecases.ListCategoriesUseCase
}

func NewCategoryFactory(db *gorm.DB, logger logging.LoggerInterface) *CategoryFactory {
	categoryRepo := repository_impl.NewCategoryRepository(db)
	return &CategoryFactory{
		CreateCategory: usecases.NewCreateCategoryUseCase(categoryRepo, logger),
		ListCategories: usecases.NewListCategoriesUseCase(categoryRepo),
	}
}
