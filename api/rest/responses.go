package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": message,
	})
}

func SendConflict(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusConflict,
		"message": message,
	})
}

func SendServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": message,
	})
}

func SendNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": message,
	})
}

func SendCreated(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusCreated, payload)
}

func SendOk(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusOK, payload)
}
