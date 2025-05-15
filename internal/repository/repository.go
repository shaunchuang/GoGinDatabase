package repository

import (
    "context"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "golang-gin-app/internal/models"
)

// Repository defines the methods for interacting with the database.
type Repository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id string) error
    BatchCreateUsers(ctx context.Context, users []*models.User) error
}

// UserRepository is the implementation of the Repository interface.
type UserRepository struct {
    db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
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

// BatchCreateUsers inserts multiple users into the database in a batch.
func (r *UserRepository) BatchCreateUsers(ctx context.Context, users []*models.User) error {
    if len(users) == 0 {
        return nil
    }
    query := `
        INSERT INTO user (account, create_time, email, last_login_date, password, status, steam_id, tel_cell, username)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    stmt, err := tx.PrepareContext(ctx, query)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()
    for _, user := range users {
        _, err := stmt.ExecContext(ctx,
            user.Account, user.CreateTime, user.Email, user.LastLoginDate,
            user.Password, user.Status, user.SteamID, user.TelCell, user.Username)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}