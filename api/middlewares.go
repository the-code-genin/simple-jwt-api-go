package api

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/services"
)

type Middlewares struct {
	usersService *services.UsersService
}

func (m *Middlewares) HandleUserAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	match := regexp.MustCompile(`^Bearer\s+([^\s]+)$`).FindStringSubmatch(authHeader)
	if len(match) != 2 {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	}

	token := match[1]
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

func NewMiddlewares(usersService *services.UsersService) *Middlewares {
	return &Middlewares{usersService}
}
