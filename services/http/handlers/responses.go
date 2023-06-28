package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type BlankStruct struct{}

func SendBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func SendConflict(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, ErrorResponse{
		Code:    http.StatusConflict,
		Message: message,
	})
}

func SendServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func SendNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func SendCreated(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusCreated, SuccessResponse{
		Code: http.StatusCreated,
		Data: payload,
	})
}

func SendOk(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusOK, SuccessResponse{
		Code: http.StatusOK,
		Data: payload,
	})
}
