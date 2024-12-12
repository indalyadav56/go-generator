package handlers

import (
	"backend/internal/auth/services"
)

type AuthHandler interface {
	// Add your handler methods here
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

// Add your handler implementations here
