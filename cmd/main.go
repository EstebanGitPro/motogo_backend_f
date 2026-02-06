package main

import (
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/server"
	"github.com/gin-gonic/gin"
)

// @title           Motogo Backend API
// @version         1.0
// @description     API RESTful para la plataforma Motogo, implementada con arquitectura hexagonal y siguiendo los principios del Richardson Maturity Model (Nivel 2-3) con HATEOAS.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Motogo API Support
// @contact.url    https://motogo.com/support
// @contact.email  support@motogo.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8085
// @BasePath  /motogo/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// CORS is configured in server.routing() to centralize all route-related configuration
	dependencies := server.Boostrap(app)

	serverAddr := dependencies.Config.GetServerAddress()
	slog.Info("Starting server", slog.String("address", serverAddr))

	if err := app.Run(serverAddr); err != nil {
		slog.Error("Server failed to start", slog.String("error", err.Error()))
		return
	}
}
