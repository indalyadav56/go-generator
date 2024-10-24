package services

import (
	"backend/internal/auth/models"
	"backend/internal/auth/repository"
	"backend/pkg/logger"
)

type AuthService interface {
	Create(auth *models.Auth) (*models.Auth, error)
	GetByID(id string) (*models.Auth, error)
	Update(id string, auth *models.Auth) (*models.Auth, error)
	Delete(id string) error
	GetAll() ([]*models.Auth, error)
	Search(query string) ([]*models.Auth, error)
}

type authService struct {
	authRepo repository.AuthRepository
	log      logger.Logger
}

func NewAuthService(repo repository.AuthRepository, log logger.Logger) *authService {
	return &authService{
		authRepo: repo,
		log:      log,
	}
}

// Create creates a new auth
func (s *authService) Create(auth *models.Auth) (*models.Auth, error) {
	return s.authRepo.Insert(auth)
}

// GetByID retrieves a auth by its ID
func (s *authService) GetByID(id string) (*models.Auth, error) {
	return s.authRepo.FindByID(id)
}

// Update updates an existing auth by its ID
func (s *authService) Update(id string, auth *models.Auth) (*models.Auth, error) {
	existingAuth, err := s.authRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields (example: only certain fields are updated)
	existingAuth.Field1 = auth.Field1
	existingAuth.Field2 = auth.Field2

	return s.authRepo.Update(id, existingAuth)
}

// Delete removes a auth by its ID
func (s *authService) Delete(id string) error {
	return s.authRepo.Delete(id)
}

// GetAll retrieves all auth records
func (s *authService) GetAll() ([]*models.Auth, error) {
	return s.authRepo.FindAll()
}

// Search allows searching for auth entities based on a query
func (s *authService) Search(query string) ([]*models.Auth, error) {
	return s.authRepo.Search(query)
}
