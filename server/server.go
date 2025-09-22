package server
import (
	"log"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/cmd/dependecy"
	"github.com/EstebanGitPro/motogo-backend/handlers/person"
	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	slog.Info("Setting up routes")

	handler := person.New(dependencies.PersonService)


	public := app.Group("/v1/motogo")
	{
		public.POST("/users", handler.RegisterPerson())
		public.GET("/users/email/:email", handler.GetPersonByEmail())
	}

	
}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		log.Fatal("Error initializing dependencies")
		return nil
	}

	routing(app, dependencies)

	return dependencies
}
