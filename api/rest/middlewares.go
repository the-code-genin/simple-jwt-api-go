package rest

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
)

type Middlewares struct {
	usersService users.UsersService
}

func (m *Middlewares) HandleUserAuth(ctx *gin.Context) {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	}

	token := strings.TrimSpace(authHeader[1])
	user, err := m.usersService.DecodeAccessToken(ctx, token)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		ctx.Abort()
		return
	}

	ctx.Set("auth_user", user)
	ctx.Set("auth_token", token)
	ctx.Next()
}

func NewMiddlewares(usersService users.UsersService) *Middlewares {
	return &Middlewares{usersService}
}
