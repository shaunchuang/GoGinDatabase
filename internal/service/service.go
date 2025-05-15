package service

import (
	"context"
	"database/sql"
	"fmt"
	"golang-gin-app/internal/models"
	"golang-gin-app/internal/repository"
	"golang-gin-app/internal/utils"
)

type Service struct {
	repo repository.Repository
	db   *sql.DB // 添加數據庫連接以便初始化計數器
}

func NewService(repo repository.Repository) *Service {
	var db *sql.DB
	db = repo.GetDB() // 通過 repository 獲取資料庫連接
	return &Service{repo: repo, db: db}
}

// GenerateFakeUsers generates a specified number of fake users and saves them to the database
func (s *Service) GenerateFakeUsers(ctx context.Context, count int, userType string, roleIDs []int64) (int, error) {
	if count < 1 || count > 1000 {
		return 0, fmt.Errorf("count must be between 1 and 1000")
	}

	// 設置預設使用者類型為 "doctor" 如果未提供
	if userType == "" {
		userType = "doctor"
	}

	// 初始化計數器，防止重複帳號
	if s.db != nil {
		if err := utils.InitializeCounters(s.db); err != nil {
			return 0, fmt.Errorf("failed to initialize counters: %v", err)
		}
	}

	// 生成假使用者
	users := utils.GenerateFakeUsers(count, userType)

	// 批量創建使用者並獲取創建後的使用者 ID
	userIDs, err := s.repo.BatchCreateUsers(ctx, users)
	if err != nil {
		return 0, fmt.Errorf("批量創建使用者失敗: %v", err)
	}

	// 如果沒有創建任何使用者，直接返回
	if len(userIDs) == 0 {
		return 0, nil
	}

	// 如果提供了角色 ID，為每個用戶分配這些角色
	if len(roleIDs) > 0 {
		// 統計角色分配失敗的用戶數
		failCount := 0
		// 直接為剛剛創建的用戶分配角色
		for _, userID := range userIDs {
			// 嘗試分配角色
			if err := s.repo.AssignRoleToUser(ctx, userID, roleIDs); err != nil {
				failCount++
				fmt.Printf("Error assigning roles to user %d: %v\n", userID, err)
				// 繼續處理下一個用戶，不中斷整個流程
			}
		}

		// 如果所有角色分配都失敗，返回錯誤
		if failCount == len(userIDs) {
			return len(users), fmt.Errorf("所有用戶角色分配失敗，請檢查 user_role 表和權限設置")
		} else if failCount > 0 {
			// 部分失敗，繼續處理但記錄警告
			fmt.Printf("警告: %d/%d 個用戶角色分配失敗\n", failCount, len(userIDs))
		}
	}

	return len(users), nil
}

// CreateUser creates a single user in the database
func (s *Service) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

// GetUserByID retrieves a user by ID from the database
func (s *Service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// ListUsers retrieves a list of users from the database
func (s *Service) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListUsersWithRoles(ctx, 50) // Limiting to 50 users for display purposes
}

// ListAllRoles 獲取所有角色
func (s *Service) ListAllRoles(ctx context.Context) ([]*models.Role, error) {
	return s.repo.ListAllRoles(ctx)
}

// AssignRolesToUser 為用戶分配角色
func (s *Service) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []int64) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleIDs)
}

// GetUserRoles 獲取用戶角色
func (s *Service) GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}
