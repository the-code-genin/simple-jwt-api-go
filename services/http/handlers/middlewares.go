package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"go.uber.org/zap"
)

type Middlewares struct {
	usersService users.UsersService
}

func (m *Middlewares) HandleUserAuth(c *gin.Context) {
	ctx := logger.With(c.Request.Context(), zap.String(logger.FunctionNameField, "Middlewares/HandleUserAuth"))

	authHeader := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		message := "invalid Authorization header"
		logger.Error(ctx, message)
		SendBadRequest(c, message)
		c.Abort()
		return
	}

	token := strings.TrimSpace(authHeader[1])
	user, err := m.usersService.DecodeAccessToken(c, token)
	if err != nil {
		message := "Unable to decode user access token"
		logger.Error(ctx, message, zap.Error(err))
		SendBadRequest(c, message)
		c.Abort()
		return
	}

	c.Set("auth_user", *user)
	c.Set("auth_token", token)
	c.Next()
}

func NewMiddlewares(usersService users.UsersService) *Middlewares {
	return &Middlewares{usersService}
}
