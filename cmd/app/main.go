package main

import (
    "github.com/gin-gonic/gin"
    "golang-gin-app/internal/app"
)

func main() {
    r := gin.Default()
    
    // Initialize the application
    app.Initialize(r)

    // Start the server
    if err := r.Run(); err != nil {
        panic(err)
    }
}