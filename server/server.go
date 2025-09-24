package server

import (
	"log"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/cmd/dependency"
	"github.com/EstebanGitPro/motogo-backend/handlers/person"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/schema"
	
	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	slog.Info("Setting up routes")

	handler := person.New(dependencies.PersonService)


	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		slog.Error("Error creating validator", slog.String("error", err.Error()))
		return
	}
	validator := middleware.NewMiddlewareValidator(validators)

	public := app.Group("/v1/motogo")
	{
		public.POST("/users", validator.WithValidateRegister(), handler.RegisterPerson())
		public.GET("/users/email/:email", handler.GetPersonByEmail())
	}

}

func Boostrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		log.Fatal("Error initializing dependencies")
		return nil
	}

	routing(app, dependencies)

	return dependencies
}
