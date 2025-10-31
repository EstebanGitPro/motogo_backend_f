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

	handler := handlers.New(dependencies.PersonService)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		slog.Error("Error creating validator", slog.String("error", err.Error()))
		log.Fatalf("Failed to initialize schema validator: %v", err)
	}
	validator := middleware.NewMiddlewareValidator(validators)

	public := app.Group("motogo/api/v1")
	{
		public.POST("/accounts", validator.WithValidateRegister(), handler.RegisterPerson())
		public.POST("/auth/login", handler.Login())
		public.GET("/accounts/email/:email", handler.GetPersonByEmail())
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
