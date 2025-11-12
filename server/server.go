package server

import (
	"log"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/cmd/dependency"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/schema"

	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	slog.Info("Setting up routes")

	app.Use(middleware.ErrorHandler())

	handler := handlers.New(dependencies.PersonService, dependencies.Interactor)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		slog.Error("Error creating validator", slog.String("error", err.Error()))
		log.Fatalf("Failed to initialize schema validator: %v", err)
	}
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

}

func Boostrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
		return nil
	}

	routing(app, dependencies)

	return dependencies
}
