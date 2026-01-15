package server

import (
	"log/slog"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/EstebanGitPro/motogo-backend/cmd/dependency"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/platform/schema"
	_ "github.com/EstebanGitPro/motogo-backend/platform/swaggo"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info(logger.LogRouteConfiguring)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8085", "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))
	dependencies.Logger.Info("CORS middleware configured")

	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.GET("/reset-password.html", func(c *gin.Context) {
		root, err := utils.FindModuleRoot()
		if err != nil {
			slog.Error("Cannot find module root", slog.String("error", err.Error()))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		resetPasswordPath := filepath.Join(root, "public", "reset-password.html")
		c.File(resetPasswordPath)
	})

	app.Use(middleware.RequestID())

	app.Use(middleware.TrackMetrics())

	errorHandler := middleware.NewErrorHandler(dependencies.MessagingCache)
	app.Use(errorHandler.Handle())

	handler := handlers.New(dependencies.Interactor, dependencies.MessageInteractor, dependencies.BranchInteractor, dependencies.BrandInteractor, dependencies.LocationInteractor, dependencies.FirebaseClient, dependencies.MessagingCache, dependencies.IDEncoder, dependencies.ResponseHandler)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		dependencies.Logger.Error(logger.LogRouteValidatorError, "error", err)
		dependencies.Logger.Fatal(logger.LogRouteValidatorError, "error", err)
	}
	dependencies.Logger.Success(logger.LogRouteValidatorOK)
	validator := middleware.NewMiddlewareValidator(validators)

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "motogo-backend",
		})
	})

	// 404 handler - must be registered AFTER all routes
	app.NoRoute(middleware.NotFoundHandler())

	// Richardson Maturity Model Nivel 2-3: Recursos con URIs únicas + HATEOAS
	public := app.Group("motogo/api/v1")
	{
		// === PERSONS ENDPOINTS ===
		// POST /persons - Crear nueva persona (registro)
		// Devuelve: 201 Created + Location header + HATEOAS links
		public.POST("/persons", validator.WithValidateRegister(), handler.RegisterPerson())

		// GET /persons/:id - Locate: Obtener persona por ID
		// Este es el endpoint referenciado en el Location header del POST
		//public.GET("/persons/:id", handler.GetPersonByID())

		// === AUTH ENDPOINTS ===
		// POST /auth/login - Autenticar usuario
		public.POST("/auth/login", handler.Login())

		// POST /auth/refresh - Refrescar access token
		public.POST("/auth/refresh", handler.RefreshToken())

		// POST /auth/resend-verification - Reenviar email de verificación
		public.POST("/auth/resend-verification", validator.WithValidateResendVerification(), handler.ResendVerificationEmail())

		// POST /auth/password-reset - Solicitar recuperación de contraseña
		public.POST("/auth/password-reset", validator.WithValidatePasswordReset(), handler.RequestPasswordReset())

		// POST /auth/verify-email - Verificar email mediante token proxy (no expone Keycloak)
		public.POST("/auth/verify-email", handler.VerifyEmailByToken())

		// POST /auth/password/reset - Actualizar contraseña con token del email de recuperación
		public.POST("/auth/password/reset", handler.ResetPasswordWithToken())

		// === MESSAGES ENDPOINTS (system administration) ===
		// POST /messages - Crear nuevo mensaje del sistema
		public.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		// PUT /messages/:id - Actualizar mensaje existente
		public.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		// DELETE /messages/:id - Eliminar mensaje
		public.DELETE("/messages/:id", handler.DeleteMessage())

		// GET /messages/:id - Obtener mensaje por ID
		public.GET("/messages/:id", handler.GetMessageByID())

		// GET /messages - Listar mensajes (con filtros opcionales)
		// Query params: ?module=users&type=ERROR&category=usuario_final&active=true
		public.GET("/messages", handler.ListMessages())

		// POST /messages/cache/reload - Recargar caché de mensajes desde BD
		// Endpoint administrativo para forzar recarga después de cambios manuales
		public.POST("/messages/cache/reload", handler.ReloadMessageCache())

		// GET /persons/:id/contact - Obtener info de contacto pública (HU55)
		public.GET("/persons/:id/contact", handler.GetPublicContact())

		// === BRANDS ENDPOINTS (catalog) ===
		// GET /brands - Listar todas las marcas disponibles
		public.GET("/brands", handler.GetBrands())

		// === GEOGRAPHIC CATALOGS (departments/cities) ===
		// GET /departments - Listar todos los departamentos
		public.GET("/departments", handler.GetDepartments())
		// GET /departments/:id/cities - Listar ciudades de un departamento
		public.GET("/departments/:id/cities", handler.GetCitiesByDepartment())

		// === DEV TOOLS (solo desarrollo) ===
		// POST /geocoding/test - Probar geocodificación sin crear sede
		public.POST("/geocoding/test", handler.TestGeocoding())

		// === BRANCH TYPES CATALOG (HU76) ===
		// GET /branch-types - Listar todos los tipos de establecimiento
		public.GET("/branch-types", handler.GetBranchTypes())
	}

	// ========================================
	// Protected Routes (require JWT authentication)
	// ========================================
	protected := app.Group("motogo/api/v1")
	protected.Use(middleware.RequireAuth(dependencies.PersonService, dependencies.MessagingCache, dependencies.JWTValidator))
	{
		// GET /persons/me - Obtener perfil del usuario autenticado
		protected.GET("/persons/me", handler.GetAuthenticatedUser())

		// PUT /persons/me - Actualizar perfil del usuario autenticado (HU52)
		protected.PUT("/persons/me", handler.UpdateProfile())

		// PUT /persons/me/password - Cambiar contraseña del usuario autenticado (HU57)
		protected.PUT("/persons/me/password", handler.ChangePassword())

		// DELETE /persons/me - Eliminar cuenta del usuario autenticado (HU53)
		protected.DELETE("/persons/me", handler.DeleteSelf())

		// GET /auth/firebase-token - Obtener token de Firebase para Storage
		protected.GET("/auth/firebase-token", handler.GetFirebaseToken())

		// === BRANCHES ENDPOINTS (HU59, HU62) ===
		// GET /branches - Listar mis sedes (solo REPRESENTANTE)
		protected.GET("/branches",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ListBranches(),
		)

		// GET /branches/:id - Consultar info de sede (HU62)
		// Accessible by all authenticated users, HATEOAS links vary by ownership
		protected.GET("/branches/:id", handler.GetBranch())

		// POST /branches - Registrar nueva sede (solo REPRESENTANTE)
		protected.POST("/branches",
			validator.WithValidateRegisterBranch(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.RegisterBranch(),
		)

		// PUT /branches/:id - Modificar sede (solo REPRESENTANTE dueño) (HU60)
		protected.PUT("/branches/:id",
			validator.WithValidateRegisterBranch(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.UpdateBranch(),
		)

		// DELETE /branches/:id - Eliminar sede (solo REPRESENTANTE dueño) (HU61)
		protected.DELETE("/branches/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeleteBranch(),
		)

		// === FRANCHISES ENDPOINTS (HU26-29) ===
		// GET /franchises - Listar mis franquicias (solo REPRESENTANTE)
		protected.GET("/franchises",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ListFranchises(dependencies.FranchiseInteractor),
		)

		// GET /franchises/:id - Consultar info de franquicia
		protected.GET("/franchises/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.GetFranchise(dependencies.FranchiseInteractor),
		)

		// POST /franchises - Registrar nueva franquicia (solo REPRESENTANTE)
		protected.POST("/franchises",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.RegisterFranchise(dependencies.FranchiseInteractor),
		)

		// PUT /franchises/:id - Modificar franquicia (solo REPRESENTANTE dueño) (HU27)
		protected.PUT("/franchises/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.UpdateFranchise(dependencies.FranchiseInteractor),
		)

		// DELETE /franchises/:id - Eliminar franquicia (solo REPRESENTANTE dueño) (HU28)
		protected.DELETE("/franchises/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeleteFranchise(dependencies.FranchiseInteractor),
		)

		// POST /franchises/:id/branches - Agregar sede a franquicia
		protected.POST("/franchises/:id/branches",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.AddBranchToFranchise(dependencies.FranchiseInteractor),
		)

		// DELETE /franchises/:id/branches/:branchId - Desvincular sede de franquicia
		protected.DELETE("/franchises/:id/branches/:branchId",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.RemoveBranchFromFranchise(dependencies.FranchiseInteractor),
		)
	}

	dependencies.Logger.Success(logger.LogRouteConfigured)
}

func Boostrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error("Failed to initialize dependencies", slog.String("error", err.Error()))
		panic(err)
	}

	routing(app, dependencies)

	return dependencies
}
