package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type BlankStruct struct{}

func SendBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, APIResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func SendPreconditionFailed(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusPreconditionFailed, APIResponse{
		Code:    http.StatusPreconditionFailed,
		Message: message,
	})
}

func SendServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, APIResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func SendNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, APIResponse{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func SendCreated(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusCreated, APIResponse{
		Code: http.StatusCreated,
		Data: payload,
	})
}

func SendOk(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusOK, APIResponse{
		Code: http.StatusOK,
		Data: payload,
	})
}
