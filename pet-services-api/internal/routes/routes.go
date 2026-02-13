package routes

import (
	"context"
	"pet-services-api/internal/config"
	"pet-services-api/internal/database"
	"pet-services-api/internal/handlers"
	"pet-services-api/internal/middlewares"
	"time"

	"pet-services-api/docs"

	"pet-services-api/internal/logging"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(storageInput database.StorageInput, ctx context.Context, logger logging.LoggerInterface) *gin.Engine {
	docs.SwaggerInfo.Title = "Pet Services API"
	docs.SwaggerInfo.Description = "API para gerenciamento de serviços pet."
	docs.SwaggerInfo.Version = "1.0"

	if host := config.GetSwaggerHost(); host != "" {
		docs.SwaggerInfo.Host = host
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
	}

	docs.SwaggerInfo.BasePath = "/"

	handlerFactory := handlers.NewHandlerFactory(storageInput, logger)
	middlewareFactory := middlewares.NewMiddlewareFactory(logger)

	r := gin.Default()

	devURL, prodURL := config.GetFrontendURLs()

	allowOrigins := []string{}
	if devURL != "" {
		allowOrigins = append(allowOrigins, devURL)
	}

	if prodURL != "" {
		allowOrigins = append(allowOrigins, prodURL)
	}

	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	allowCredentials := len(allowOrigins) > 0 && allowOrigins[0] != "*"

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewareFactory.DefaultRateLimitMiddleware())

	public := r.Group("/")
	{
		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		public.GET("/health", handlerFactory.HealthHandler.HealthCheck)

		public.POST("/users/register", handlerFactory.UserHandler.RegisterUser)
		public.POST("/users/check-email", handlerFactory.UserHandler.CheckEmailExists)
		public.POST("/users/check-phone", handlerFactory.UserHandler.CheckPhoneExists)
	}

	authPublic := r.Group("/auth/")
	authPublic.Use(middlewareFactory.StrictRateLimitMiddleware())
	{
		authPublic.POST("login", handlerFactory.TokenHandler.LoginUser)
		authPublic.POST("refresh", handlerFactory.TokenHandler.RefreshToken)
		authPublic.POST("request-password-reset", handlerFactory.TokenHandler.RequestPasswordReset)
		authPublic.POST("reset-password", handlerFactory.TokenHandler.ResetPassword)
		authPublic.POST("resend-verification-email", handlerFactory.TokenHandler.ResendVerificationEmail)
		authPublic.POST("verify-email", handlerFactory.TokenHandler.VerifyEmail)
	}

	authorizedAuth := r.Group("/auth/")
	authorizedAuth.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedAuth.POST("logout", handlerFactory.TokenHandler.Logout)
	}

	authorizedUser := r.Group("/users/")
	authorizedUser.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedUser.GET("/profile", handlerFactory.UserHandler.GetProfile)
		authorizedUser.GET("/:user_id", handlerFactory.UserHandler.GetUserByID)
		authorizedUser.GET("", handlerFactory.UserHandler.ListUsers)
		authorizedUser.PUT("", handlerFactory.UserHandler.UpdateUser)
		authorizedUser.DELETE("", handlerFactory.UserHandler.DeleteUser)
		authorizedUser.POST("/reactivate", handlerFactory.UserHandler.ReactivateUser)
		authorizedUser.POST("/deactivate", handlerFactory.UserHandler.DeactivateUser)
		authorizedUser.POST("/change-password", handlerFactory.UserHandler.ChangePassword)
		authorizedUser.POST("/update-email-verified", handlerFactory.UserHandler.UpdateEmailVerified)
		authorizedUser.POST("/photos", handlerFactory.UserHandler.AddUserPhoto)
	}

	authorizedReference := r.Group("/reference/")
	authorizedReference.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedReference.GET("/countries", handlerFactory.ReferenceHandler.ListCountries)
		authorizedReference.GET("/states", handlerFactory.ReferenceHandler.ListStates)
		authorizedReference.GET("/cities", handlerFactory.ReferenceHandler.ListCities)
	}

	authorizedSpecies := r.Group("/species/")
	authorizedSpecies.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedSpecies.GET("", handlerFactory.SpecieHandler.ListSpecies)
	}

	authorizedOwner := r.Group("/pets/")
	authorizedOwner.Use(middlewareFactory.AuthMiddleware(), middlewareFactory.OwnerOnlyMiddleware())
	{
		authorizedOwner.POST("", handlerFactory.PetHandler.AddPet)
		authorizedOwner.POST("/:pet_id/photos", handlerFactory.PetHandler.AddPetPhoto)
	}

	authorizedOwnerProviders := r.Group("/providers/")
	authorizedOwnerProviders.Use(middlewareFactory.AuthMiddleware(), middlewareFactory.OwnerOnlyMiddleware())
	{
		authorizedOwnerProviders.POST("/:provider_id/reviews", handlerFactory.ReviewHandler.CreateReview)
	}

	authorizedProvider := r.Group("/providers/")
	authorizedProvider.Use(middlewareFactory.AuthMiddleware(), middlewareFactory.ProviderOnlyMiddleware())
	{
		authorizedProvider.POST("", handlerFactory.ProviderHandler.AddProvider)
		authorizedProvider.POST("/photos", handlerFactory.ProviderHandler.AddProviderPhoto)
	}

	authorizedServices := r.Group("/services/")
	authorizedServices.Use(middlewareFactory.AuthMiddleware(), middlewareFactory.ProviderOnlyMiddleware())
	{
		authorizedServices.POST("", handlerFactory.ServiceHandler.AddService)
		authorizedServices.POST("/:service_id/photos", handlerFactory.ServiceHandler.AddServicePhoto)
		authorizedServices.POST("/:service_id/tags", handlerFactory.ServiceHandler.AddServiceTag)
		authorizedServices.POST("/:service_id/categories", handlerFactory.ServiceHandler.AddServiceCategory)
	}

	r.GET("/tags", middlewareFactory.AuthMiddleware(), middlewareFactory.ProviderOnlyMiddleware(), handlerFactory.ServiceHandler.ListTags)
	r.GET("/categories", middlewareFactory.AuthMiddleware(), middlewareFactory.ProviderOnlyMiddleware(), handlerFactory.CategoryHandler.ListCategories)

	authorizedRequests := r.Group("/requests/")
	authorizedRequests.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedRequests.GET("", handlerFactory.RequestHandler.ListRequests)
		authorizedRequests.POST("", middlewareFactory.OwnerOnlyMiddleware(), handlerFactory.RequestHandler.AddRequest)
	}

	authorizedAdmin := r.Group("/admin/")
	authorizedAdmin.Use(middlewareFactory.AuthMiddleware(), middlewareFactory.AdminOnlyMiddleware())
	{
		authorizedAdmin.POST("", handlerFactory.UserHandler.CreateAdmin)
		authorizedAdmin.POST("/categories", handlerFactory.CategoryHandler.CreateCategory)
	}

	return r
}
