package handlers

import (
	"backend/internal/user/services"
)

type UserHandler interface {
	// Add your handler methods here
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

// Add your handler implementations here
