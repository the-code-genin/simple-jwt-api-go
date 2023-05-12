package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
)

type UsersAuthHandlers struct {
	usersService users.UsersService
}

func (a *UsersAuthHandlers) HandleRegister(ctx *gin.Context) {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "UsersAuthHandlers/HandleRegister")

	var req users.RegisterUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error(err.Error())
		SendBadRequest(ctx, err.Error())
		return
	}

	log = log.WithField(logger.RequestBodyField, req)

	user, err := a.usersService.Register(ctx, req)
	if err != nil {
		log.WithError(err).Error(err.Error())
		switch err {
		case users.ErrEmailTaken:
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
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "UsersAuthHandlers/HandleGenerateAccessToken")

	var req users.GenerateUserAccessTokenDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error(err.Error())
		SendBadRequest(ctx, err.Error())
		return
	}

	user, token, err := a.usersService.GenerateAccessToken(ctx, req)
	if err != nil {
		log.WithError(err).Error(err.Error())
		switch err {
		case users.ErrInvalidPassword:
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
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "UsersAuthHandlers/HandleBlacklistAccessToken")

	val, ok := ctx.Get("auth_token")
	if !ok {
		log.Error("Auth token not in gin context")
		SendServerError(ctx, "an error occured")
		return
	}

	token, ok := val.(string)
	if !ok {
		log.Error("Unable to parse auth token from gin context")
		SendServerError(ctx, "an error occured")
		return
	}

	err := a.usersService.BlacklistAccessToken(ctx, token)
	if err != nil {
		log.WithError(err).Error(err.Error())
		SendServerError(ctx, err.Error())
		return
	}

	SendOk(ctx, gin.H{})
}

func (a *UsersAuthHandlers) HandleGetMe(ctx *gin.Context) {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "UsersAuthHandlers/HandleGetMe")

	val, ok := ctx.Get("auth_user")
	if !ok {
		log.Error("Auth token not in gin context")
		SendServerError(ctx, "an error occured")
		return
	}

	authUser, ok := val.(*users.UserDTO)
	if !ok {
		log.Error("Unable to parse auth user from gin context")
		SendServerError(ctx, "an error occured")
		return
	}

	SendOk(ctx, gin.H{
		"user": authUser,
	})
}

func NewUsersAuthHandlers(
	usersService users.UsersService,
) *UsersAuthHandlers {
	return &UsersAuthHandlers{usersService}
}
