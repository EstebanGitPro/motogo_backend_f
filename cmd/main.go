package main

import (
    "log/slog"
    "github.com/EstebanGitPro/motogo-backend/server"
    "github.com/gin-gonic/gin"
)

func main() {
    gin.SetMode(gin.ReleaseMode)

    app := gin.New()
    app.Use(gin.Logger())
    app.Use(gin.Recovery())

    dependencies := server.Boostrap(app)

    serverAddr := dependencies.Config.GetServerAddress()
    slog.Info("Starting server", slog.String("address", serverAddr))

    if err := app.Run(serverAddr); err != nil {
        slog.Error("Server failed to start", slog.String("error", err.Error()))
        return
    }
}