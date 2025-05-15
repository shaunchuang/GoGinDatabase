package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// HelloHandler handles the /hello route
func HelloHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

// GetUserHandler handles the /user/:id route
func GetUserHandler(c *gin.Context) {
    id := c.Param("id")
    // Here you would typically fetch the user from the database
    c.JSON(http.StatusOK, gin.H{"id": id, "name": "John Doe"})
}

// CreateUserHandler handles the POST /user route
func CreateUserHandler(c *gin.Context) {
    var user struct {
        Name string `json:"name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Here you would typically save the user to the database
    c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
}