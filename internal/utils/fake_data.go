package utils

import (
	"golang-gin-app/internal/models"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// GenerateFakeUser creates a fake user with realistic data
func GenerateFakeUser() *models.User {
	now := time.Now()
	pastDate := gofakeit.DateRange(now.AddDate(-5, 0, 0), now)
	steamID := gofakeit.UUID()
	return &models.User{
		Account:       gofakeit.Username(),
		CreateTime:    now,
		Email:         gofakeit.Email(),
		LastLoginDate: &pastDate,
		Password:      gofakeit.Password(true, true, true, false, false, 12),
		Status:        "active",
		SteamID:       &steamID,
		TelCell:       gofakeit.Phone(),
		Username:      gofakeit.Name(),
	}
}

// GenerateFakeUsers creates a slice of fake users
func GenerateFakeUsers(count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = GenerateFakeUser()
	}
	return users
}
