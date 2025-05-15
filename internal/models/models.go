package models

import "time"

type User struct {
    ID            int64     `json:"id"`
    Account       string    `json:"account"`
    CreateTime    time.Time `json:"create_time"`
    Email         string    `json:"email"`
    LastLoginDate *time.Time `json:"last_login_date"`
    Password      string    `json:"password"`
    Status        string    `json:"status"`
    SteamID       *string   `json:"steam_id"`
    TelCell       string    `json:"tel_cell"`
    Username      string    `json:"username"`
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