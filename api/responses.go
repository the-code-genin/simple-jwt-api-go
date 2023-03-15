package api

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

func SendServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": message,
	})
}
