package handlers

import (
	"golang-gin-app/internal/service"
	"golang-gin-app/internal/utils"
	"net/http"
	"strconv"
	"strings"

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

// GenerateFakeUsersFormHandler handles the GET /fake-users route to display the form
func GenerateFakeUsersFormHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取所有角色
		roles, err := svc.ListAllRoles(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Failed to fetch roles: " + err.Error(),
			})
			return
		}

		// 獲取用戶列表
		users, err := svc.ListUsers(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Failed to fetch user list: " + err.Error(),
				"roles": roles,
			})
			return
		}
		c.HTML(http.StatusOK, "fake_users.html", gin.H{
			"title": "Generate Fake Users",
			"users": users,
			"roles": roles,
		})
	}
}

// GenerateFakeUsersHandler handles the POST /fake-users route to generate fake users
func GenerateFakeUsersHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		countStr := c.PostForm("count")
		userType := c.PostForm("userType") // 取得使用者類型

		// 取得選擇的角色 ID
		roleIDsStr := c.PostFormArray("roleIDs")
		roleIDs := make([]int64, 0, len(roleIDsStr))

		for _, idStr := range roleIDsStr {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err == nil { // 忽略無效的 ID
				roleIDs = append(roleIDs, id)
			}
		}

		// 如果沒有選擇角色，則使用預設角色
		if len(roleIDs) == 0 {
			roleIDs = utils.GetDefaultRoleIDsForUserType(userType)
		}

		// 如果沒有選擇角色，則使用預設角色
		if len(roleIDs) == 0 {
			roleIDs = utils.GetDefaultRoleIDsForUserType(userType)
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 {
			c.HTML(http.StatusBadRequest, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Please enter a valid number greater than 0",
			})
			return
		}

		// 將使用者類型和角色傳遞給 service 方法
		createdCount, err := svc.GenerateFakeUsers(c.Request.Context(), count, userType, roleIDs)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Failed to generate users: " + err.Error(),
			})
			return
		}

		// 獲取所有角色供表單使用
		roles, err := svc.ListAllRoles(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Failed to fetch roles: " + err.Error(),
			})
			return
		}

		users, err := svc.ListUsers(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_users.html", gin.H{
				"title": "Generate Fake Users",
				"error": "Failed to fetch user list after generation: " + err.Error(),
				"roles": roles,
			})
			return
		} // 生成角色名稱列表，用於訊息顯示
		roleNames := make([]string, 0, len(roleIDs))
		roleIDMap := make(map[int64]string)

		for _, role := range roles {
			roleIDMap[role.ID] = role.Alias
		}

		for _, id := range roleIDs {
			if alias, ok := roleIDMap[id]; ok {
				roleNames = append(roleNames, alias)
			}
		}

		// 顯示已生成使用者的類型、數量和角色
		message := strconv.Itoa(createdCount) + " fake " + userType + " users"
		if len(roleNames) > 0 {
			message += " with roles: " + strings.Join(roleNames, ", ")
		}
		message += " have been successfully generated and saved to the database."

		c.HTML(http.StatusOK, "fake_users.html", gin.H{
			"title":   "Generate Fake Users",
			"message": message,
			"users":   users,
			"roles":   roles,
		})
	}
}
