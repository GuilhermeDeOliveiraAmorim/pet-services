package routes

import (
	"context"
	"pet-services-api/internal/config"
	"pet-services-api/internal/database"
	"pet-services-api/internal/handlers"
	"pet-services-api/internal/middlewares"
	"pet-services-api/internal/repository_impl"
	"strings"
	"time"

	"pet-services-api/docs"

	"pet-services-api/internal/logging"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func SetupRouter(storageInput database.StorageInput, ctx context.Context, logger logging.LoggerInterface) *gin.Engine {
	docs.SwaggerInfo.Title = "Pet Services API"
	docs.SwaggerInfo.Description = "API para gerenciamento de serviços pet."
	docs.SwaggerInfo.Version = "1.0"

	if host := config.GetSwaggerHost(); host != "" {
		docs.SwaggerInfo.Host = host
		if strings.Contains(host, "onrender.com") {
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else {
			docs.SwaggerInfo.Schemes = []string{"http"}
		}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	docs.SwaggerInfo.BasePath = "/"

	handlerFactory := handlers.NewHandlerFactory(storageInput, logger)
	middlewareFactory := middlewares.NewMiddlewareFactory(logger)
	userRepo := repository_impl.NewUserRepository(storageInput.DB)
	profileComplete := middlewares.ProfileCompleteMiddleware(logger, userRepo)

	r := gin.Default()

	// Desabilitar redirecionamento automático de trailing slash para evitar problemas com CORS
	r.RedirectTrailingSlash = false

	corsOrigins := config.GetCORSOrigins()
	allowOrigins := []string{}

	if corsOrigins != "" {
		// Split CORS_ORIGINS by comma
		for _, origin := range splitAndTrim(corsOrigins, ",") {
			if origin != "" {
				allowOrigins = append(allowOrigins, origin)
			}
		}
	} else {
		// Fallback to old env vars
		devURL, prodURL := config.GetFrontendURLs()
		if devURL != "" {
			allowOrigins = append(allowOrigins, devURL)
		}
		if prodURL != "" {
			allowOrigins = append(allowOrigins, prodURL)
		}
	}

	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	allowCredentials := len(allowOrigins) > 0 && allowOrigins[0] != "*"

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
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

	authPublic := r.Group("/auth")
	authPublic.Use(middlewareFactory.StrictRateLimitMiddleware())
	{
		authPublic.POST("login", handlerFactory.TokenHandler.LoginUser)
		authPublic.POST("refresh", handlerFactory.TokenHandler.RefreshToken)
		authPublic.POST("request-password-reset", handlerFactory.TokenHandler.RequestPasswordReset)
		authPublic.POST("reset-password", handlerFactory.TokenHandler.ResetPassword)
		authPublic.POST("resend-verification-email", handlerFactory.TokenHandler.ResendVerificationEmail)
		authPublic.POST("verify-email", handlerFactory.TokenHandler.VerifyEmail)
	}

	authorizedAuth := r.Group("/auth")
	authorizedAuth.Use(middlewareFactory.AuthMiddleware())
	{
		authorizedAuth.POST("logout", handlerFactory.TokenHandler.Logout)
	}

	authorizedUser := r.Group("/users")
	authorizedUser.Use(middlewareFactory.AuthMiddleware(), profileComplete)
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

	publicReference := r.Group("/reference")
	{
		publicReference.GET("/countries", handlerFactory.ReferenceHandler.ListCountries)
		publicReference.GET("/states", handlerFactory.ReferenceHandler.ListStates)
		publicReference.GET("/cities", handlerFactory.ReferenceHandler.ListCities)
	}

	// Rotas utilitárias públicas - sem autenticação
	publicUtil := r.Group("/util")
	{
		publicUtil.GET("/species", handlerFactory.SpecieHandler.ListSpecies)
		publicUtil.GET("/categories", handlerFactory.CategoryHandler.ListCategories)
	}

	authorizedOwner := r.Group("/pets")
	authorizedOwner.Use(middlewareFactory.AuthMiddleware(), profileComplete, middlewareFactory.OwnerOnlyMiddleware())
	{
		authorizedOwner.GET("", handlerFactory.PetHandler.ListPets)
		authorizedOwner.GET("/:pet_id", handlerFactory.PetHandler.GetPet)
		authorizedOwner.PUT("/:pet_id", handlerFactory.PetHandler.UpdatePet)
		authorizedOwner.DELETE("/:pet_id", handlerFactory.PetHandler.DeletePet)
		authorizedOwner.POST("", handlerFactory.PetHandler.AddPet)
		authorizedOwner.POST("/:pet_id/photos", handlerFactory.PetHandler.AddPetPhoto)
		authorizedOwner.DELETE("/:pet_id/photos/:photo_id", handlerFactory.PetHandler.DeletePetPhoto)
	}

	authorizedOwnerProviders := r.Group("/providers")
	authorizedOwnerProviders.Use(middlewareFactory.AuthMiddleware(), profileComplete, middlewareFactory.OwnerOnlyMiddleware())
	{
		authorizedOwnerProviders.POST("/:provider_id/reviews", handlerFactory.ReviewHandler.CreateReview)
	}

	authorizedProvider := r.Group("/providers")
	authorizedProvider.Use(middlewareFactory.AuthMiddleware(), profileComplete, middlewareFactory.ProviderOnlyMiddleware())
	{
		authorizedProvider.POST("", handlerFactory.ProviderHandler.AddProvider)
		authorizedProvider.POST("/photos", handlerFactory.ProviderHandler.AddProviderPhoto)
		authorizedProvider.DELETE("/:provider_id/photos/:photo_id", handlerFactory.ProviderHandler.DeleteProviderPhoto)
		authorizedProvider.DELETE("/:provider_id", handlerFactory.ProviderHandler.DeleteProvider)
	}

	authorizedServices := r.Group("/services")
	authorizedServices.Use(middlewareFactory.AuthMiddleware(), profileComplete, middlewareFactory.ProviderOnlyMiddleware())
	{
		authorizedServices.POST("", handlerFactory.ServiceHandler.AddService)
		authorizedServices.PUT("/:service_id", handlerFactory.ServiceHandler.UpdateService)
		authorizedServices.DELETE("/:service_id", handlerFactory.ServiceHandler.DeleteService)
		authorizedServices.POST("/:service_id/photos", handlerFactory.ServiceHandler.AddServicePhoto)
		authorizedServices.DELETE("/:service_id/photos/:photo_id", handlerFactory.ServiceHandler.DeleteServicePhoto)
		authorizedServices.POST("/:service_id/tags", handlerFactory.ServiceHandler.AddServiceTag)
		authorizedServices.POST("/:service_id/categories", handlerFactory.ServiceHandler.AddServiceCategory)
		authorizedServices.DELETE("/:service_id/categories/:category_id", handlerFactory.ServiceHandler.RemoveServiceCategory)
	}

	r.GET("/services", handlerFactory.ServiceHandler.ListServices)
	r.GET("/services/:service_id", handlerFactory.ServiceHandler.GetService)
	r.GET("/services/search", handlerFactory.ServiceHandler.SearchServices)
	r.GET("/providers/:provider_id", handlerFactory.ProviderHandler.GetProvider)
	r.GET("/reviews", handlerFactory.ReviewHandler.ListReviews)
	r.GET("/tags", handlerFactory.ServiceHandler.ListTags)

	authorizedRequests := r.Group("/requests")
	authorizedRequests.Use(middlewareFactory.AuthMiddleware(), profileComplete)
	{
		authorizedRequests.GET("", handlerFactory.RequestHandler.ListRequests)
		authorizedRequests.GET("/:request_id", handlerFactory.RequestHandler.GetRequest)
		authorizedRequests.POST("", middlewareFactory.OwnerOnlyMiddleware(), handlerFactory.RequestHandler.AddRequest)
		authorizedRequests.PATCH("/:request_id/accept", middlewareFactory.ProviderOnlyMiddleware(), handlerFactory.RequestHandler.AcceptRequest)
		authorizedRequests.PATCH("/:request_id/reject", middlewareFactory.ProviderOnlyMiddleware(), handlerFactory.RequestHandler.RejectRequest)
		authorizedRequests.PATCH("/:request_id/complete", middlewareFactory.ProviderOnlyMiddleware(), handlerFactory.RequestHandler.CompleteRequest)
	}

	authorizedAdmin := r.Group("/admin")
	authorizedAdmin.Use(middlewareFactory.AuthMiddleware(), profileComplete, middlewareFactory.AdminOnlyMiddleware())
	{
		authorizedAdmin.POST("", handlerFactory.UserHandler.CreateAdmin)
		authorizedAdmin.POST("/categories", handlerFactory.CategoryHandler.CreateCategory)
	}

	return r
}
