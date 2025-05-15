package service

import (
	"context"
	"fmt"
	"golang-gin-app/internal/models"
	"golang-gin-app/internal/repository"
	"golang-gin-app/internal/utils"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// GenerateFakeUsers generates a specified number of fake users and saves them to the database
func (s *Service) GenerateFakeUsers(ctx context.Context, count int) (int, error) {
	if count < 1 || count > 1000 {
		return 0, fmt.Errorf("count must be between 1 and 1000")
	}
	users := utils.GenerateFakeUsers(count)
	err := s.repo.BatchCreateUsers(ctx, users)
	if err != nil {
		return 0, err
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
