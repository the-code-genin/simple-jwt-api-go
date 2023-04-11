package api

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database/repositories"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	config *internal.Config
	users  *repositories.Users
	tokens *repositories.BlacklistedTokens
}

func (a *AuthHandlers) HandleSignup(ctx *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	// Check if the email is taken
	emailTaken, err := a.users.EmailTaken(body.Email)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		return
	} else if emailTaken {
		SendBadRequest(ctx, "Email not available.")
		return
	}

	// Hash the user's password
	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	// Insert the user record
	user, err := a.users.Insert(&repositories.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: hex.EncodeToString(password),
	})
	if err != nil {
		SendServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func (a *AuthHandlers) HandleLogin(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	// Get the user and verify the password
	user, err := a.users.GetUserWithEmail(body.Email)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}
	hashedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}
	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(body.Password)); err != nil {
		SendBadRequest(ctx, "Invalid password.")
		return
	}

	// Generate JWT token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Second * time.Duration(a.config.JWT.Exp)).Unix(),
	}).SignedString(a.config.JWT.Key)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":         user,
		"access_token": token,
		"type":         "bearer",
	})
}

func (a *AuthHandlers) HandleGetMe(ctx *gin.Context) {
	val, ok := ctx.Get("auth_user")
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}
	authUser, ok := val.(*repositories.User)
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": authUser,
	})
}

func (a *AuthHandlers) HandleLogout(ctx *gin.Context) {
	val, ok := ctx.Get("auth_token")
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}
	authToken, ok := val.(*jwt.Token)
	if !ok {
		SendServerError(ctx, "an error occured")
		return
	}

	claims := authToken.Claims.(jwt.MapClaims)
	err := a.tokens.Add(authToken.Raw, int64(claims["exp"].(float64)))
	if err != nil {
		SendServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func NewAuthHandlers(
	config *internal.Config,
	users *repositories.Users,
	tokens *repositories.BlacklistedTokens,
) *AuthHandlers {
	return &AuthHandlers{config, users, tokens}
}
