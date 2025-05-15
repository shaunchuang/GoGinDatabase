package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "golang-gin-app/internal/service"
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

// GenerateFakeUsersFormHandler handles the GET /fake-users route to display the form
func GenerateFakeUsersFormHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "fake_users.html", gin.H{
        "title": "Generate Fake Users",
    })
}

// GenerateFakeUsersHandler handles the POST /fake-users route to generate fake users
func GenerateFakeUsersHandler(svc *service.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        countStr := c.PostForm("count")
        count, err := strconv.Atoi(countStr)
        if err != nil || count < 1 {
            c.HTML(http.StatusBadRequest, "fake_users.html", gin.H{
                "title": "Generate Fake Users",
                "error": "Please enter a valid number greater than 0",
            })
            return
        }
        createdCount, err := svc.GenerateFakeUsers(c.Request.Context(), count)
        if err != nil {
            c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
                "title": "Generate Fake Users",
                "error": "Failed to generate users: " + err.Error(),
            })
            return
        }
        c.HTML(http.StatusOK, "fake_users.html", gin.H{
            "title":   "Generate Fake Users",
            "message": strconv.Itoa(createdCount) + " fake users have been successfully generated and saved to the database.",
        })
    }
}