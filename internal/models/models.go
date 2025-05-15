package models

import "time"

type User struct {
	ID            int64      `json:"id"`
	Account       string     `json:"account"`
	CreateTime    time.Time  `json:"create_time"`
	Email         string     `json:"email"`
	LastLoginDate *time.Time `json:"last_login_date"`
	Password      string     `json:"password"`
	Status        string     `json:"status"`
	SteamID       *string    `json:"steam_id"`
	TelCell       string     `json:"tel_cell"`
	Username      string     `json:"username"`
	// 新增欄位，用於前端顯示角色
	Roles []*Role `json:"roles,omitempty"`
}

type Role struct {
	ID          int64   `json:"id"`
	Alias       string  `json:"alias"`
	Description *string `json:"description,omitempty"`
}

type UserRole struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	RoleID    int64      `json:"role_id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}
