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

	handler := handlers.New(dependencies.PersonService, dependencies.Interactor, dependencies.MessageInteractor, dependencies.MessagingCache, dependencies.IDEncoder, dependencies.ResponseHandler)

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
	}

	// ========================================
	// Protected Routes (require JWT authentication)
	// ========================================
	protected := app.Group("motogo/api/v1")
	protected.Use(middleware.RequireAuth(dependencies.PersonService, dependencies.MessagingCache))
	{
		// GET /persons/me - Obtener perfil del usuario autenticado (alias contextual)
		protected.GET("/persons/me", handler.GetAuthenticatedUser())
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
