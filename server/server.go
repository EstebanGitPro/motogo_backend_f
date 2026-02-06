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

	handler := handlers.New(
		dependencies.Interactor,
		dependencies.MessageInteractor,
		dependencies.BranchInteractor,
		dependencies.BrandInteractor,
		dependencies.LocationInteractor,
		dependencies.ServiceInteractor,
		dependencies.FranchiseInteractor,
		dependencies.MotorcycleInteractor,
		dependencies.EvidenceInteractor, // HU16-19
		dependencies.FirebaseClient,
		dependencies.MessagingCache,
		dependencies.IDEncoder,
		dependencies.ResponseHandler,
	)

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

	// 404 handler - must be registered AFTER all rote
	app.NoRoute(middleware.NotFoundHandler())

	public := app.Group("motogo/api/v1")
	{
		// === PERSONS ENDPOINTS ===
		// POST /persons - Crear nueva persona (registro)
		// Devuelve: 201 Created + Location header + HATEOAS links
		public.POST("/persons", validator.WithValidateRegister(), handler.RegisterPerson())

		// GET /persons/:id - Locate: Obtener persona por ID
		// Este es el endpoint referenciado en el Location header del POST
		// public.GET("/persons/:id", handler.GetPersonByID())

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
		// Endpoint administration para forzar recarga después de cambios manuals
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

		// === BRANCH TYPES CATALOG (HU76) ===
		// GET /branch-types - Listar todos los tipos de establecimiento
		public.GET("/branch-types", handler.GetBranchTypes())

		// === SERVICES CATALOG (HU63, HU75) ===
		// GET /service-types - Listar todos los tipos de servicio (HU75)
		public.GET("/service-types", handler.GetServiceTypes())
		// GET /services - Listar catálogo de servicios con filtro opcional (HU63)
		public.GET("/services", handler.GetServices())
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

		// === GEOCODING (any authenticated user) ===
		// POST /location/geocode - Geocodificar dirección para mapas
		protected.POST("/location/geocode", handler.TestGeocoding())

		// === BRANCHES ENDPOINTS (HU59, HU62) ===
		// GET /branches - Listar mis sedes (solo REPRESENTANTE)
		protected.GET("/branches",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ListBranches(),
		)

		// GET /branches/nearby - Buscar sedes cercanas (HU89)
		// Accessible by all authenticated users (USER, REPRESENTATIVE)
		protected.GET("/branches/nearby", handler.GetNearbyBranches())

		// GET /branches/:id - Consultar info de sede (HU62)
		// Accessible by all authenticated users, HATEOAS links vary by ownership
		protected.GET("/branches/:id", handler.GetBranch())

		// GET /branches/:id/services - Obtener servicios de una sede
		// Accessible by all authenticated users
		protected.GET("/branches/:id/services", handler.GetBranchServices())

		// === MOTORCYCLES ENDPOINTS (HU43-47) ===
		// Only USER role can access motorcycle endpoints (not ADMIN, not REPRESENTATIVE)

		// POST /motorcycles - Register new motorcycle (HU43)
		protected.POST("/motorcycles",
			middleware.RequireRole(domain.RoleUser),
			validator.WithValidateRegisterMotorcycle(),
			handler.RegisterMotorcycle(),
		)

		// GET /motorcycles/:id - Get motorcycle details (HU46)
		// HATEOAS links vary by ownership (Richardson Level 3)
		protected.GET("/motorcycles/:id",
			middleware.RequireRole(domain.RoleUser),
			handler.GetMotorcycle(),
		)

		// GET /motorcycles - List all motorcycles for authenticated user (HU46)
		protected.GET("/motorcycles",
			middleware.RequireRole(domain.RoleUser),
			handler.ListMotorcycles(),
		)

		// PUT /motorcycles/:id - Update motorcycle (HU44, solo propietario)
		protected.PUT("/motorcycles/:id",
			middleware.RequireRole(domain.RoleUser),
			handler.UpdateMotorcycle(),
		)

		// DELETE /motorcycles/:id - Soft delete motorcycle (HU45, solo propietario)
		protected.DELETE("/motorcycles/:id",
			middleware.RequireRole(domain.RoleUser),
			handler.DeleteMotorcycle(),
		)

		// === MOTORCYCLE PROFILE IMAGE ENDPOINTS (HU36-39) ===
		// Only USER role can manage their motorcycle's profile image

		// PUT /motorcycles/:id/profile-image - Add or update profile image (HU36/37)
		protected.PUT("/motorcycles/:id/profile-image",
			middleware.RequireRole(domain.RoleUser),
			handler.UpdateProfileImage(),
		)

		// GET /motorcycles/:id/profile-image - Get profile image URL (HU38)
		protected.GET("/motorcycles/:id/profile-image",
			middleware.RequireRole(domain.RoleUser),
			handler.GetProfileImage(),
		)

		// DELETE /motorcycles/:id/profile-image - Remove profile image (HU39)
		protected.DELETE("/motorcycles/:id/profile-image",
			middleware.RequireRole(domain.RoleUser),
			handler.DeleteProfileImage(),
		)

		// GET /motorcycles/lookup?plate={placa} - Lookup motorcycle by plate (HU47)
		// Accessible by representatives (workshops) for service purposes
		protected.GET("/motorcycles/lookup",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.LookupMotorcycleByPlate(),
		)

		// GET /motorcycle-references - List motorcycle reference catalog (HU50)
		// Accessible by USER role for motorcycle registration
		protected.GET("/motorcycle-references",
			middleware.RequireRole(domain.RoleUser),
			handler.GetMotorcycleReferences(),
		)

		// === MOTORCYCLE EVIDENCE ENDPOINTS (HU16-19) ===
		// Photographic evidence gallery for motorcycles

		// POST /motorcycles/:id/evidence - Upload evidence for a motorcycle (HU16)
		protected.POST("/motorcycles/:id/evidence",
			middleware.RequireRole(domain.RoleUser),
			validator.WithValidateEvidence(),
			handler.CreateEvidence(),
		)

		// GET /motorcycles/:id/evidence - List evidence for a motorcycle (HU18)
		protected.GET("/motorcycles/:id/evidence",
			middleware.RequireRole(domain.RoleUser),
			handler.ListEvidence(),
		)

		// DELETE /motorcycles/:id/evidence/:evidenceId - Delete evidence (HU19)
		protected.DELETE("/motorcycles/:id/evidence/:evidenceId",
			middleware.RequireRole(domain.RoleUser),
			handler.DeleteEvidence(),
		)

		// POST /branches/:id/services - Asociar servicios a una sede (solo REPRESENTANTE)
		protected.POST("/branches/:id/services",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.AssociateBranchServices(),
		)

		// DELETE /branches/:id/services/:serviceId - Desasociar servicio de una sede (solo REPRESENTANTE)
		protected.DELETE("/branches/:id/services/:serviceId",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DissociateBranchService(),
		)

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

		// === SCHEDULES ENDPOINTS (HU30-35) ===
		// POST /branches/:id/schedules - Crear horario para una sede (HU30)
		protected.POST("/branches/:id/schedules",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.CreateBranchSchedule(dependencies.ScheduleInteractor),
		)

		// GET /branches/:id/schedules - Consultar horario de una sede (HU32)
		protected.GET("/branches/:id/schedules",
			handler.GetBranchSchedule(dependencies.ScheduleInteractor),
		)

		// PUT /branches/:id/schedules - Modificar horario de una sede (HU31)
		protected.PUT("/branches/:id/schedules",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.UpdateBranchSchedule(dependencies.ScheduleInteractor),
		)

		// DELETE /branches/:id/schedules - Eliminar horario de una sede (HU33)
		protected.DELETE("/branches/:id/schedules",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeleteBranchSchedule(dependencies.ScheduleInteractor),
		)

		// PUT /branches/:id/schedules/activate - Activar horario (HU34)
		protected.PUT("/branches/:id/schedules/activate",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ActivateBranchSchedule(dependencies.ScheduleInteractor),
		)

		// PUT /branches/:id/schedules/deactivate - Desactivar horario (HU35)
		protected.PUT("/branches/:id/schedules/deactivate",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeactivateBranchSchedule(dependencies.ScheduleInteractor),
		)

		// === SCHEDULE DETAILS ENDPOINTS (HU6-9) ===
		// POST /branches/:id/schedules/details - Crear franja horaria (HU6)
		protected.POST("/branches/:id/schedules/details",
			validator.WithValidateScheduleDetail(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.CreateScheduleDetail(dependencies.ScheduleDetailInteractor, dependencies.ScheduleInteractor),
		)

		// GET /branches/:id/schedules/details - Listar franjas horarias (HU9)
		// Accessible by all authenticated users (USER can view schedule, REPRESENTATIVE can manage)
		protected.GET("/branches/:id/schedules/details",
			handler.ListScheduleDetails(dependencies.ScheduleDetailInteractor, dependencies.ScheduleInteractor),
		)

		// PUT /schedule-details/:id - Modificar franja horaria (HU7)
		protected.PUT("/schedule-details/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.UpdateScheduleDetail(dependencies.ScheduleDetailInteractor),
		)

		// DELETE /schedule-details/:id - Eliminar franja horaria (HU8)
		protected.DELETE("/schedule-details/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeleteScheduleDetail(dependencies.ScheduleDetailInteractor),
		)

		// GET /schedules/days - Catálogo de días de la semana (HU10)
		protected.GET("/schedules/days", handler.GetDaysOfWeek())

		// === SCHEDULE EXCEPTIONS ENDPOINTS (HU20-25) ===
		// POST /branches/:id/schedules/exceptions - Crear excepción de horario (HU20)
		protected.POST("/branches/:id/schedules/exceptions",
			validator.WithValidateScheduleException(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.CreateScheduleException(dependencies.ScheduleExceptionInteractor, dependencies.ScheduleInteractor),
		)

		// GET /branches/:id/schedules/exceptions - Listar excepciones de horario (HU23)
		protected.GET("/branches/:id/schedules/exceptions",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ListScheduleExceptions(dependencies.ScheduleExceptionInteractor, dependencies.ScheduleInteractor),
		)

		// PUT /schedule-exceptions/:id - Modificar excepción de horario (HU21)
		protected.PUT("/schedule-exceptions/:id",
			validator.WithValidateUpdateScheduleException(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.UpdateScheduleException(dependencies.ScheduleExceptionInteractor),
		)

		// DELETE /schedule-exceptions/:id - Eliminar excepción de horario (HU22)
		protected.DELETE("/schedule-exceptions/:id",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeleteScheduleException(dependencies.ScheduleExceptionInteractor),
		)

		// PUT /schedule-exceptions/:id/activate - Activar excepción de horario (HU24)
		protected.PUT("/schedule-exceptions/:id/activate",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.ActivateScheduleException(dependencies.ScheduleExceptionInteractor),
		)

		// PUT /schedule-exceptions/:id/deactivate - Desactivar excepción de horario (HU25)
		protected.PUT("/schedule-exceptions/:id/deactivate",
			middleware.RequireRole(domain.RoleRepresentative),
			handler.DeactivateScheduleException(dependencies.ScheduleExceptionInteractor),
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
			validator.WithValidateFranchise(),
			middleware.RequireRole(domain.RoleRepresentative),
			handler.RegisterFranchise(dependencies.FranchiseInteractor),
		)

		// PUT /franchises/:id - Modificar franquicia (solo REPRESENTANTE dueño) (HU27)
		protected.PUT("/franchises/:id",
			validator.WithValidateFranchise(),
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

	// ========================================
	// Admin Routes (require JWT + ADMIN role)
	// ========================================
	admin := app.Group("motogo/api/v1/admin")
	admin.Use(middleware.RequireAuth(dependencies.PersonService, dependencies.MessagingCache, dependencies.JWTValidator))
	admin.Use(middleware.RequireRole(domain.RoleAdmin))
	{
		// === SERVICES CATALOG ADMIN (HU68) ===
		// PUT /admin/services/:id - Modificar servicio del catálogo global
		admin.PUT("/services/:id", handler.UpdateService())

		// === BRAND LINES ADMIN (HU40) ===
		// GET /admin/brands/:brandId/lines - Consultar líneas de una marca
		admin.GET("/brands/:brandId/lines", handler.GetBrandLines())
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
