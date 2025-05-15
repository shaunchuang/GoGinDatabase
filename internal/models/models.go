package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

type Post struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    UserID  int    `json:"user_id"`
}

type Comment struct {
    ID     int    `json:"id"`
    PostID int    `json:"post_id"`
    UserID int    `json:"user_id"`
    Content string `json:"content"`
}