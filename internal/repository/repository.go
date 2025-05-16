package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-gin-app/internal/models"
	"time"

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

	// 可預約時段相關方法
	BatchCreateAvailableSlots(ctx context.Context, slots []*models.AvailableSlot) error
	GetAvailableSlotsByDoctor(ctx context.Context, doctorID int64) ([]*models.AvailableSlot, error)
	GetUserByRoleID(ctx context.Context, roleID int64) ([]*models.User, error)
	UpdateAvailableSlot(ctx context.Context, slot *models.AvailableSlot) error
	DeleteAvailableSlot(ctx context.Context, slotID int64) error
	GetAvailableSlotByID(ctx context.Context, slotID int64) (*models.AvailableSlot, error)
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

// GetUserByRoleID 獲取特定角色的用戶列表
func (r *UserRepository) GetUserByRoleID(ctx context.Context, roleID int64) ([]*models.User, error) {
	query := `
		SELECT u.ID, u.account, u.create_time, u.email, u.last_login_date, 
		       u.password, u.status, u.steam_id, u.tel_cell, u.username
		FROM user u
		JOIN user_role ur ON u.ID = ur.user_id
		WHERE ur.role_id = ?
		ORDER BY u.ID DESC
	`
	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("獲取角色用戶列表失敗: %v", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Account, &user.CreateTime, &user.Email, &user.LastLoginDate,
			&user.Password, &user.Status, &user.SteamID, &user.TelCell, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("掃描用戶數據失敗: %v", err)
		}

		// 獲取該用戶的角色
		roles, err := r.GetUserRoles(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("獲取用戶 %d 的角色失敗: %v", user.ID, err)
		}
		user.Roles = roles

		users = append(users, user)
	}
	return users, nil
}

// BatchCreateAvailableSlots 批量創建可預約時段
func (r *UserRepository) BatchCreateAvailableSlots(ctx context.Context, slots []*models.AvailableSlot) error {
	if len(slots) == 0 {
		return nil
	}

	// 開始事務
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開始事務失敗: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			fmt.Printf("事務回滾，錯誤: %v\n", err)
		}
	}()

	// 檢查表是否存在
	checkTableQuery := `SHOW TABLES LIKE 'wg_available_slots'`
	rows, err := r.db.QueryContext(ctx, checkTableQuery)
	if err != nil {
		return fmt.Errorf("檢查 wg_available_slots 表失敗: %v", err)
	}

	tableExists := false
	for rows.Next() {
		tableExists = true
		break
	}
	rows.Close()

	if !tableExists {
		return fmt.Errorf("wg_available_slots 表不存在，請創建表")
	}

	// 批量插入時段
	insertQuery := `INSERT INTO wg_available_slots 
                   (doctor, is_booked, slot_begin_time, slot_date, slot_end_time) 
                   VALUES (?, ?, ?, ?, ?)`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("準備插入語句失敗: %v", err)
	}
	defer stmt.Close()

	for _, slot := range slots {
		_, err = stmt.ExecContext(ctx,
			slot.Doctor,
			slot.IsBooked,
			slot.SlotBeginTime.Format("15:04:05"),
			slot.SlotDate.Format("2006-01-02"),
			slot.SlotEndTime.Format("15:04:05"))
		if err != nil {
			return fmt.Errorf("插入時段失敗: %v", err)
		}
	}

	// 提交事務
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事務失敗: %v", err)
	}

	return nil
}

// GetAvailableSlotsByDoctor 獲取醫師的可預約時段
func (r *UserRepository) GetAvailableSlotsByDoctor(ctx context.Context, doctorID int64) ([]*models.AvailableSlot, error) {
	query := `
		SELECT ID, doctor, is_booked, slot_begin_time, slot_date, slot_end_time
		FROM wg_available_slots
		WHERE doctor = ?
		ORDER BY slot_date, slot_begin_time
	`
	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("獲取醫師時段失敗: %v", err)
	}
	defer rows.Close()

	slots := make([]*models.AvailableSlot, 0)
	for rows.Next() {
		slot := &models.AvailableSlot{}
		var beginTime, endTime, slotDate string
		err := rows.Scan(
			&slot.ID,
			&slot.Doctor,
			&slot.IsBooked,
			&beginTime,
			&slotDate,
			&endTime)
		if err != nil {
			return nil, fmt.Errorf("掃描時段數據失敗: %v", err)
		}

		// 解析日期時間 - 嘗試多種格式
		// 首先嘗試標準日期格式
		date, err := time.Parse("2006-01-02", slotDate)
		if err != nil {
			// 如果標準格式解析失敗，嘗試 ISO 格式
			date, err = time.Parse(time.RFC3339, slotDate)
			if err != nil {
				// 嘗試不帶時區的格式
				date, err = time.Parse("2006-01-02T15:04:05", slotDate)
				if err != nil {
					return nil, fmt.Errorf("解析日期失敗 %s: %v", slotDate, err)
				}
			}
		}
		// 只保留日期部分，去除時間
		slot.SlotDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

		// 解析開始時間
		beginTimeParsed, err := time.Parse("15:04:05", beginTime)
		if err != nil {
			return nil, fmt.Errorf("解析開始時間失敗 %s: %v", beginTime, err)
		}
		slot.SlotBeginTime = time.Date(
			date.Year(), date.Month(), date.Day(),
			beginTimeParsed.Hour(), beginTimeParsed.Minute(), beginTimeParsed.Second(),
			0, time.Local)

		// 解析結束時間
		endTimeParsed, err := time.Parse("15:04:05", endTime)
		if err != nil {
			return nil, fmt.Errorf("解析結束時間失敗 %s: %v", endTime, err)
		}
		slot.SlotEndTime = time.Date(
			date.Year(), date.Month(), date.Day(),
			endTimeParsed.Hour(), endTimeParsed.Minute(), endTimeParsed.Second(),
			0, time.Local)

		slots = append(slots, slot)
	}
	return slots, nil
}

// UpdateAvailableSlot 更新可預約時段
func (r *UserRepository) UpdateAvailableSlot(ctx context.Context, slot *models.AvailableSlot) error {
	query := `
		UPDATE wg_available_slots
		SET doctor = ?, is_booked = ?, slot_begin_time = ?, slot_date = ?, slot_end_time = ?
		WHERE ID = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		slot.Doctor,
		slot.IsBooked,
		slot.SlotBeginTime.Format("15:04:05"),
		slot.SlotDate.Format("2006-01-02"),
		slot.SlotEndTime.Format("15:04:05"),
		slot.ID)

	if err != nil {
		return fmt.Errorf("更新時段失敗: %v", err)
	}
	return nil
}

// DeleteAvailableSlot 刪除可預約時段
func (r *UserRepository) DeleteAvailableSlot(ctx context.Context, slotID int64) error {
	query := `DELETE FROM wg_available_slots WHERE ID = ?`
	result, err := r.db.ExecContext(ctx, query, slotID)
	if err != nil {
		return fmt.Errorf("刪除時段失敗: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("獲取影響行數失敗: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("未找到ID為 %d 的時段", slotID)
	}

	return nil
}

// GetAvailableSlotByID 通過ID獲取時段
func (r *UserRepository) GetAvailableSlotByID(ctx context.Context, slotID int64) (*models.AvailableSlot, error) {
	query := `
		SELECT ID, doctor, is_booked, slot_begin_time, slot_date, slot_end_time
		FROM wg_available_slots
		WHERE ID = ?
	`
	slot := &models.AvailableSlot{}
	var beginTime, endTime, slotDate string

	err := r.db.QueryRowContext(ctx, query, slotID).Scan(
		&slot.ID,
		&slot.Doctor,
		&slot.IsBooked,
		&beginTime,
		&slotDate,
		&endTime)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("未找到ID為 %d 的時段", slotID)
	}

	if err != nil {
		return nil, fmt.Errorf("獲取時段數據失敗: %v", err)
	}

	// 解析日期時間 - 嘗試多種格式
	// 首先嘗試標準日期格式
	date, err := time.Parse("2006-01-02", slotDate)
	if err != nil {
		// 如果標準格式解析失敗，嘗試 ISO 格式
		date, err = time.Parse(time.RFC3339, slotDate)
		if err != nil {
			// 嘗試不帶時區的格式
			date, err = time.Parse("2006-01-02T15:04:05", slotDate)
			if err != nil {
				return nil, fmt.Errorf("解析日期失敗 %s: %v", slotDate, err)
			}
		}
	}
	// 只保留日期部分，去除時間
	slot.SlotDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	// 解析開始時間
	beginTimeParsed, err := time.Parse("15:04:05", beginTime)
	if err != nil {
		return nil, fmt.Errorf("解析開始時間失敗 %s: %v", beginTime, err)
	}
	slot.SlotBeginTime = time.Date(
		date.Year(), date.Month(), date.Day(),
		beginTimeParsed.Hour(), beginTimeParsed.Minute(), beginTimeParsed.Second(),
		0, time.Local)

	// 解析結束時間
	endTimeParsed, err := time.Parse("15:04:05", endTime)
	if err != nil {
		return nil, fmt.Errorf("解析結束時間失敗 %s: %v", endTime, err)
	}
	slot.SlotEndTime = time.Date(
		date.Year(), date.Month(), date.Day(),
		endTimeParsed.Hour(), endTimeParsed.Minute(), endTimeParsed.Second(),
		0, time.Local)

	return slot, nil
}
