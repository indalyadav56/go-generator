package controllers

import (
	"backend/internal/auth/services"
	"backend/pkg/logger"
	"backend/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type authController struct {
	log     logger.Logger
	authSrv services.AuthService
}

func NewAuthController(log logger.Logger, s services.AuthService) *authController {
	return &authController{
		authSrv: s,
		log:     log,
	}
}

// Create godoc
// @Summary		Create a resource
// @Description	Create a new auth
// @Tags			Auth
// @Produce		json
// @Router			/v1/auths [post]
// @Success		200	"Resource created successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *authController) Create(ctx *gin.Context) {
	// c.authSrv.Create()
	resp := response.Created("Successfully Created", nil)
	ctx.JSON(resp.Status, resp)
}

// Get godoc
// @Summary		Get a resource
// @Description	Fetch a auth
// @Tags			Auth
// @Produce		json
// @Router			/v1/auths [get]
// @Success		200	"Resource fetched successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *authController) Get(ctx *gin.Context) {
	resp := response.Success("Successfully fetched data", nil)
	ctx.JSON(resp.Status, resp)
}

// Update godoc
// @Summary		Update a resource
// @Description	Update a auth
// @Tags			Auth
// @Produce		json
// @Router			/v1/auths [patch]
// @Success		200	"Resource updated successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *authController) Update(ctx *gin.Context) {
	resp := response.Success("Successfully updated data", nil)
	ctx.JSON(resp.Status, resp)
}

// Delete godoc
// @Summary		Delete a resource
// @Description	Delete a auth
// @Tags			Auth
// @Produce		json
// @Router			/v1/auths [delete]
// @Success		200	"Resource deleted successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *authController) Delete(ctx *gin.Context) {
	resp := response.Success("Successfully deleted data", nil)
	ctx.JSON(resp.Status, resp)
}
