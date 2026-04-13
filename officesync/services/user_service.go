package services

import (
	"errors"
	"strings"

	"officesync/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(name, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	if name == "" || email == "" {
		return nil, errors.New("name and email are required fields")
	}

	user := &models.User{
		Name:  name,
		Email: email,
		Role:  "employee",
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "unique constraint") || strings.Contains(result.Error.Error(), "duplicate key") {
			return nil, errors.New("a user with this email already exists")
		}
		return nil, errors.New("failed to create user in the database")
	}

	return user, nil
}
