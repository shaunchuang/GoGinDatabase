package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-gin-app/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Repository defines the methods for interacting with the database.
type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	BatchCreateUsers(ctx context.Context, users []*models.User) ([]int64, error)
	ListUsers(ctx context.Context, limit int) ([]*models.User, error)
	// 新增角色相關方法
	ListAllRoles(ctx context.Context) ([]*models.Role, error)
	AssignRoleToUser(ctx context.Context, userID int64, roleIDs []int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error)
	ListUsersWithRoles(ctx context.Context, limit int) ([]*models.User, error)
	GetDB() *sql.DB // 新增方法以獲取資料庫連接
}

// UserRepository is the implementation of the Repository interface.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetDB returns the database connection
func (r *UserRepository) GetDB() *sql.DB {
	return r.db
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO user (account, create_time, email, last_login_date, password, status, steam_id, tel_cell, username)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		user.Account, user.CreateTime, user.Email, user.LastLoginDate,
		user.Password, user.Status, user.SteamID, user.TelCell, user.Username)
	return err
}

// GetByID retrieves a user by their ID from the database.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT ID, account, create_time, email, last_login_date, password, status, steam_id, tel_cell, username
              FROM user WHERE ID = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Account, &user.CreateTime, &user.Email, &user.LastLoginDate,
		&user.Password, &user.Status, &user.SteamID, &user.TelCell, &user.Username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// Update modifies an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE user SET account = ?, create_time = ?, email = ?, last_login_date = ?,
        password = ?, status = ?, steam_id = ?, tel_cell = ?, username = ?
        WHERE ID = ?`
	_, err := r.db.ExecContext(ctx, query,
		user.Account, user.CreateTime, user.Email, user.LastLoginDate,
		user.Password, user.Status, user.SteamID, user.TelCell, user.Username, user.ID)
	return err
}

// Delete removes a user from the database by their ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user WHERE ID = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// BatchCreateUsers inserts multiple users into the database in a batch and returns their IDs.
func (r *UserRepository) BatchCreateUsers(ctx context.Context, users []*models.User) ([]int64, error) {
	if len(users) == 0 {
		return []int64{}, nil
	}
	query := `
        INSERT INTO user (account, create_time, email, last_login_date, password, status, steam_id, tel_cell, username)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		res, err := stmt.ExecContext(ctx,
			user.Account, user.CreateTime, user.Email, user.LastLoginDate,
			user.Password, user.Status, user.SteamID, user.TelCell, user.Username)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 獲取新插入用戶的 ID
		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return userIDs, nil
}

// ListUsers retrieves a list of users from the database, limited by the specified number.
func (r *UserRepository) ListUsers(ctx context.Context, limit int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit if none specified
	}
	query := `SELECT ID, account, create_time, email, last_login_date, password, status, steam_id, tel_cell, username
              FROM user ORDER BY ID DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Account, &user.CreateTime, &user.Email, &user.LastLoginDate,
			&user.Password, &user.Status, &user.SteamID, &user.TelCell, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// ListAllRoles 獲取所有角色
func (r *UserRepository) ListAllRoles(ctx context.Context) ([]*models.Role, error) {
	query := `SELECT ID, alias, description FROM role`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*models.Role, 0)
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Alias, &role.Description)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// AssignRoleToUser 為用戶指派角色
func (r *UserRepository) AssignRoleToUser(ctx context.Context, userID int64, roleIDs []int64) error {
	// 檢查用戶角色表是否存在
	checkTableQuery := `SHOW TABLES LIKE 'user_role'`
	rows, err := r.db.QueryContext(ctx, checkTableQuery)
	if err != nil {
		return fmt.Errorf("檢查 user_role 表失敗: %v", err)
	}

	tableExists := false
	for rows.Next() {
		tableExists = true
		break
	}
	rows.Close()

	if !tableExists {
		return fmt.Errorf("user_role 表不存在，請創建表")
	}

	// 打印詳細信息
	fmt.Printf("正在為用戶 ID %d 分配角色: %v\n", userID, roleIDs)

	// 開始事務
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開始事務失敗: %v", err)
	}
	defer func() {
		if err != nil {
			// 發生錯誤時回滾事務
			tx.Rollback()
			fmt.Printf("事務回滾，用戶 ID: %d, 錯誤: %v\n", userID, err)
		}
	}()

	// 先刪除該用戶現有的所有角色
	deleteQuery := `DELETE FROM user_role WHERE user_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		return fmt.Errorf("刪除用戶舊角色失敗: %v", err)
	}

	// 如果沒有要添加的角色，直接提交事務
	if len(roleIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("提交空角色事務失敗: %v", err)
		}
		return nil
	}

	// 添加新角色 - 根據實際表結構修改 (移除 created_at 欄位)
	insertQuery := `INSERT INTO user_role (user_id, role_id) VALUES (?, ?)`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("準備插入語句失敗: %v", err)
	}
	defer stmt.Close()

	for _, roleID := range roleIDs {
		_, err = stmt.ExecContext(ctx, userID, roleID)
		if err != nil {
			return fmt.Errorf("為用戶 %d 插入角色 %d 失敗: %v", userID, roleID, err)
		}
		fmt.Printf("成功為用戶 %d 添加角色 %d\n", userID, roleID)
	}

	// 提交事務
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事務失敗: %v", err)
	}

	fmt.Printf("成功完成用戶 %d 的角色分配\n", userID)
	return nil
}

// GetUserRoles 獲取用戶角色
func (r *UserRepository) GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error) {
	query := `
		SELECT r.ID, r.alias, r.description
		FROM role r
		JOIN user_role ur ON r.ID = ur.role_id
		WHERE ur.user_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*models.Role, 0)
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Alias, &role.Description)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// ListUsersWithRoles 獲取包含角色資訊的用戶列表
func (r *UserRepository) ListUsersWithRoles(ctx context.Context, limit int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit if none specified
	}

	// 先獲取用戶基本資訊
	users, err := r.ListUsers(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 為每個用戶添加角色資訊
	for _, user := range users {
		roles, err := r.GetUserRoles(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
	}

	return users, nil
}
