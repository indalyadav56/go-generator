package controllers

import (
	"backend/internal/todo/services"
	"backend/pkg/logger"
	"backend/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type TodoController interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type todoController struct {
	log     logger.Logger
	todoSrv services.TodoService
}

func NewTodoController(log logger.Logger, s services.TodoService) *todoController {
	return &todoController{
		todoSrv: s,
		log:     log,
	}
}

// Create godoc
// @Summary		Create a resource
// @Description	Create a new todo
// @Tags			Todo
// @Produce		json
// @Router			/v1/todos [post]
// @Success		200	"Resource created successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *todoController) Create(ctx *gin.Context) {
	// c.todoSrv.Create()
	resp := response.Created("Successfully Created", nil)
	ctx.JSON(resp.Status, resp)
}

// Get godoc
// @Summary		Get a resource
// @Description	Fetch a todo
// @Tags			Todo
// @Produce		json
// @Router			/v1/todos [get]
// @Success		200	"Resource fetched successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *todoController) Get(ctx *gin.Context) {
	resp := response.Success("Successfully fetched data", nil)
	ctx.JSON(resp.Status, resp)
}

// Update godoc
// @Summary		Update a resource
// @Description	Update a todo
// @Tags			Todo
// @Produce		json
// @Router			/v1/todos [patch]
// @Success		200	"Resource updated successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *todoController) Update(ctx *gin.Context) {
	resp := response.Success("Successfully updated data", nil)
	ctx.JSON(resp.Status, resp)
}

// Delete godoc
// @Summary		Delete a resource
// @Description	Delete a todo
// @Tags			Todo
// @Produce		json
// @Router			/v1/todos [delete]
// @Success		200	"Resource deleted successfully"
// @Failure		400	"Bad request"
// @Failure		500	"Internal server error"
func (c *todoController) Delete(ctx *gin.Context) {
	resp := response.Success("Successfully deleted data", nil)
	ctx.JSON(resp.Status, resp)
}
