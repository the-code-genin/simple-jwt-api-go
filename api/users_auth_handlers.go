package api

import (
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/database/users"
	"github.com/the-code-genin/simple-jwt-api-go/services"
)

type UsersAuthHandlers struct {
	usersService *services.UsersService
}

func (a *UsersAuthHandlers) HandleRegister(ctx *gin.Context) {
	var req services.RegisterUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	user, err := a.usersService.Register(ctx, req)
	if err != nil {
		switch err {
		case services.ErrEmailTaken:
			SendConflict(ctx, err.Error())
		default:
			SendServerError(ctx, err.Error())
		}
		return
	}

	SendCreated(ctx, gin.H{
		"user": user,
	})
}

func (a *UsersAuthHandlers) HandleGenerateAccessToken(ctx *gin.Context) {
	var req services.GenerateUserAccessTokenDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	user, token, err := a.usersService.GenerateAccessToken(ctx, req)
	if err != nil {
		switch err {
		case services.ErrInvalidPassword:
			SendBadRequest(ctx, err.Error())
		default:
			SendServerError(ctx, err.Error())
		}
		return
	}

	SendOk(ctx, gin.H{
		"user":         user,
		"access_token": token,
		"type":         "bearer",
	})
}

func (a *UsersAuthHandlers) HandleBlacklistAccessToken(ctx *gin.Context) {
	val, ok := ctx.Get("auth_token")
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	token, ok := val.(string)
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	err := a.usersService.BlacklistAccessToken(ctx, token)
	if err != nil {
		SendServerError(ctx, err.Error())
		return
	}

	SendOk(ctx, gin.H{})
}

func (a *UsersAuthHandlers) HandleGetMe(ctx *gin.Context) {
	val, ok := ctx.Get("auth_user")
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	authUser, ok := val.(*users.User)
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	SendOk(ctx, gin.H{
		"user": authUser,
	})
}

func NewUsersAuthHandlers(
	usersService *services.UsersService,
) *UsersAuthHandlers {
	return &UsersAuthHandlers{usersService}
}
