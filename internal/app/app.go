package app

import (
    "github.com/gin-gonic/gin"
    "golang-gin-app/configs"
    "golang-gin-app/pkg/middleware"
)

type App struct {
    Router *gin.Engine
}

func NewApp() *App {
    router := gin.Default()
    app := &App{Router: router}

    app.initializeRoutes()
    app.initializeMiddleware()

    return app
}

func (a *App) initializeRoutes() {
    // Initialize your routes here
}

func (a *App) initializeMiddleware() {
    // Initialize your middleware here
    a.Router.Use(middleware.SomeMiddleware())
}

func (a *App) Run(addr string) error {
    return a.Router.Run(addr)
}