package service

import (
    "golang-gin-app/internal/models"
    "golang-gin-app/internal/repository"
)

type Service struct {
    repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
    return &Service{repo: repo}
}

// Example business logic function
func (s *Service) CreateItem(item models.Item) error {
    return s.repo.Create(item)
}

// Example business logic function
func (s *Service) GetItem(id string) (models.Item, error) {
    return s.repo.GetByID(id)
}

// Additional business logic functions can be added here.