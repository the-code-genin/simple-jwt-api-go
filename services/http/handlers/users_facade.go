package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"go.uber.org/zap"
)

type UsersFacade struct {
	usersService users.UsersService
}

// Register godoc
//
// @Summary Register a new user
// @Accept  json
// @Produce json
// @Param   req body      users.RegisterUserDTO true "body"
// @Success 200 {object} APIResponse{data=users.UserDTO}
// @Failure 400 {object} APIResponse
// @Failure 409 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router  /register  [post]
func (a *UsersFacade) Register(c *gin.Context) {
	ctx := logger.With(c.Request.Context(), zap.String(logger.FunctionNameField, "UsersFacade/Register"))

	var req users.RegisterUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Unable to bind request body to users.RegisterUserDTO", zap.Error(err))
		SendBadRequest(c, err.Error())
		return
	}

	user, err := a.usersService.Register(c, req)
	if err != nil {
		message := "An error occured while registering the user"
		logger.Error(ctx, message, zap.Error(err))
		SendPreconditionFailed(c, err.Error())
		return
	}

	SendCreated(c, user)
}

// GenerateAccessToken godoc
//
// @Summary Generate access token for a new user
// @Accept  json
// @Produce json
// @Param   req body      users.GenerateUserAccessTokenDTO true "body"
// @Success 200 {object} APIResponse{data=BlankStruct{user=users.UserDTO,access_token=string,type=string}}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router  /generate-access-token  [post]
func (a *UsersFacade) GenerateAccessToken(c *gin.Context) {
	ctx := logger.With(c.Request.Context(), zap.String(logger.FunctionNameField, "UsersFacade/GenerateAccessToken"))

	var req users.GenerateUserAccessTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Unable to bind request body to users.GenerateUserAccessTokenDTO", zap.Error(err))
		SendBadRequest(c, err.Error())
		return
	}

	user, token, err := a.usersService.GenerateAccessToken(c, req)
	if err != nil {
		message := "An error occured while generate user access token"
		logger.Error(ctx, message, zap.Error(err))
		SendPreconditionFailed(c, err.Error())
		return
	}

	SendOk(c, gin.H{
		"user":         user,
		"access_token": token,
		"type":         "bearer",
	})
}

// BlacklistAccessToken godoc
//
// @Summary  Blacklist access token for user
// @Produce  json
// @security securitydefinitions.apikey
// @Success  200 {object} APIResponse{data=BlankStruct}
// @Failure  400 {object} APIResponse
// @Failure  500 {object} APIResponse
// @Router   /blacklist-access-token  [post]
func (a *UsersFacade) BlacklistAccessToken(c *gin.Context) {
	ctx := logger.With(c.Request.Context(), zap.String(logger.FunctionNameField, "UsersFacade/BlacklistAccessToken"))

	val, ok := c.Get("auth_token")
	if !ok {
		logger.Error(ctx, "Auth token not in gin context")
		SendServerError(c, "an error occured")
		return
	}

	token, ok := val.(string)
	if !ok {
		logger.Error(ctx, "Unable to parse auth token from gin context")
		SendServerError(c, "an error occured")
		return
	}

	err := a.usersService.BlacklistAccessToken(c, token)
	if err != nil {
		message := "An error occured while blacklisting user access token"
		logger.Error(ctx, message, zap.Error(err))
		SendPreconditionFailed(c, err.Error())
		return
	}

	SendOk(c, BlankStruct{})
}

// GetMe godoc
//
// @Summary  Get authenticated user
// @Produce  json
// @security securitydefinitions.apikey
// @Success  200 {object} APIResponse{data=users.UserDTO}
// @Failure  400 {object} APIResponse
// @Failure  500 {object} APIResponse
// @Router   /me [get]
func (a *UsersFacade) GetMe(c *gin.Context) {
	ctx := logger.With(c.Request.Context(), zap.String(logger.FunctionNameField, "UsersFacade/GetMe"))

	val, ok := c.Get("auth_user")
	if !ok {
		logger.Error(ctx, "Auth token not in gin context")
		SendServerError(c, "an error occured")
		return
	}

	authUser, ok := val.(users.UserDTO)
	if !ok {
		logger.Error(ctx, "Unable to parse auth user from gin context")
		SendServerError(c, "an error occured")
		return
	}

	SendOk(c, authUser)
}

func NewUsersFacade(
	usersService users.UsersService,
) *UsersFacade {
	return &UsersFacade{usersService}
}
