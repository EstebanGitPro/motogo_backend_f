package server

import (
	"github.com/EstebanGitPro/motogo-backend/cmd/dependency"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/schema"

	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info("Configurando rutas de la aplicación")

	app.Use(middleware.ErrorHandler(dependencies.Logger))

	handler := handlers.New(dependencies.PersonService, dependencies.Interactor, dependencies.Logger, dependencies.IDEncoder)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		dependencies.Logger.Error("Error creando validador de schema", "error", err)
		dependencies.Logger.Fatal("Failed to initialize schema validator", "error", err)
	}
	dependencies.Logger.Success("Validador de schema inicializado")
	validator := middleware.NewMiddlewareValidator(validators)

	// Richardson Maturity Model Nivel 2-3: Recursos con URIs únicas + HATEOAS
	public := app.Group("motogo/api/v1")
	{
		// POST /accounts - Crear nueva cuenta
		// Devuelve: 201 Created + Location header + HATEOAS links
		public.POST("/accounts", validator.WithValidateRegister(), handler.RegisterPerson())

		// GET /accounts/:id - Locate: Obtener cuenta por ID
		// Este es el endpoint referenciado en el Location header del POST
		//public.GET("/accounts/:id", handler.GetPersonByID())

		//public.POST("/auth/login", handler.Login())
		//public.GET("/accounts/email/:email", handler.GetPersonByEmail())
	}

	dependencies.Logger.Success("Rutas configuradas correctamente")
}

func Boostrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		dependencies.Logger.Fatal("Error initializing dependencies", "error", err)
		return nil
	}

	routing(app, dependencies)

	return dependencies
}
