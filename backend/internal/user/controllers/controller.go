package controllers

import (
	"backend/internal/user/services"
	"backend/pkg/logger"
	"backend/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type userController struct {
	log     logger.Logger
	userSrv services.UserService
}

func NewUserController(log logger.Logger, s services.UserService) *userController {
	return &userController{
		userSrv: s,
		log:     log,
	}
}

// Create godoc
// @Summary		Create a resource
// @Description	Create a new user
// @Tags			User
// @Produce		json
// @Router			/v1/users [post]
// @Success		200	"Resource created successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *userController) Create(ctx *gin.Context) {
	// c.userSrv.Create()
	resp := response.Created("Successfully Created", nil)
	ctx.JSON(resp.Status, resp)
}

// Get godoc
// @Summary		Get a resource
// @Description	Fetch a user
// @Tags			User
// @Produce		json
// @Router			/v1/users [get]
// @Success		200	"Resource fetched successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *userController) Get(ctx *gin.Context) {
	resp := response.Success("Successfully fetched data", nil)
	ctx.JSON(resp.Status, resp)
}

// Update godoc
// @Summary		Update a resource
// @Description	Update a user
// @Tags			User
// @Produce		json
// @Router			/v1/users [patch]
// @Success		200	"Resource updated successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *userController) Update(ctx *gin.Context) {
	resp := response.Success("Successfully updated data", nil)
	ctx.JSON(resp.Status, resp)
}

// Delete godoc
// @Summary		Delete a resource
// @Description	Delete a user
// @Tags			User
// @Produce		json
// @Router			/v1/users [delete]
// @Success		200	"Resource deleted successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *userController) Delete(ctx *gin.Context) {
	resp := response.Success("Successfully deleted data", nil)
	ctx.JSON(resp.Status, resp)
}
