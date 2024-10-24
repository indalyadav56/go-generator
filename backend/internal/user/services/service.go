package services

import (
	"backend/internal/user/models"
	"backend/internal/user/repository"
	"backend/pkg/logger"
)

type UserService interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) (*models.User, error)
	Delete(id string) error
	GetAll() ([]*models.User, error)
	Search(query string) ([]*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	log      logger.Logger
}

func NewUserService(repo repository.UserRepository, log logger.Logger) *userService {
	return &userService{
		userRepo: repo,
		log:      log,
	}
}

// Create creates a new user
func (s *userService) Create(user *models.User) (*models.User, error) {
	return s.userRepo.Insert(user)
}

// GetByID retrieves a user by its ID
func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// Update updates an existing user by its ID
func (s *userService) Update(id string, user *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields (example: only certain fields are updated)
	existingUser.Field1 = user.Field1
	existingUser.Field2 = user.Field2

	return s.userRepo.Update(id, existingUser)
}

// Delete removes a user by its ID
func (s *userService) Delete(id string) error {
	return s.userRepo.Delete(id)
}

// GetAll retrieves all user records
func (s *userService) GetAll() ([]*models.User, error) {
	return s.userRepo.FindAll()
}

// Search allows searching for user entities based on a query
func (s *userService) Search(query string) ([]*models.User, error) {
	return s.userRepo.Search(query)
}
