package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleSuccess(ctx *gin.Context, message string, data ...interface{}) {
	var payload interface{}
	if len(data) > 0 {
		payload = data[0]
	}

	resp := Success(message, payload)
	ctx.JSON(http.StatusOK, resp)
}

func HandleCreated(ctx *gin.Context, message string, data ...interface{}) {
	var payload interface{}
	if len(data) > 0 {
		payload = data[0]
	}

	resp := Success(message, payload)
	ctx.JSON(http.StatusCreated, resp)
}

func HandleError(ctx *gin.Context, statusCode int, message string, errors ...interface{}) {
	var payload interface{}
	if len(errors) > 0 {
		payload = errors[0]
	}

	resp := Error(statusCode, message, payload)
	ctx.JSON(resp.Status, resp)
}

func HandleErrorWithAbort(ctx *gin.Context, statusCode int, message string, errors ...interface{}) {
	var payload interface{}
	if len(errors) > 0 {
		payload = errors[0]
	}

	resp := Error(statusCode, message, payload)
	ctx.JSON(resp.Status, resp)
	ctx.Abort()
}

func HandleUnAuthorizedErrorWithAbort(ctx *gin.Context, message string) {
	resp := Unauthorized(message)
	ctx.JSON(resp.Status, resp)
	ctx.Abort()
}

func HandleServerErrorWithAbort(ctx *gin.Context, message string) {
	resp := Unauthorized(message)
	ctx.JSON(resp.Status, resp)
	ctx.Abort()
}
