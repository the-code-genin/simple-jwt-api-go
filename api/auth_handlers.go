package api

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
	"golang.org/x/crypto/bcrypt"
)

func NewSignupHandler(users *database.Users) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		emailTaken, err := users.EmailTaken(body.Email)
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
		user := &database.User{}
		user.Name = body.Name
		user.Email = body.Email
		user.Password = hex.EncodeToString(password)
		user, err = users.Insert(user)
		if err != nil {
			SendServerError(ctx, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"user": user,
		})
	}
}

func NewLoginHandler(config *internal.Config, users *database.Users) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			SendBadRequest(ctx, err.Error())
			return
		}

		// Get the user
		user, err := users.GetUserWithEmail(body.Email)
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
			SendBadRequest(ctx, err.Error())
			return
		}

		// Generate JWT token
		jwtKey, err := config.GetJWTKey()
		if err != nil {
			SendBadRequest(ctx, err.Error())
			return
		}
		jwtExpiry, err := config.GetJWTExpiry()
		if err != nil {
			SendBadRequest(ctx, err.Error())
			return
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":    user.ID,
			"user_email": user.Email,
			"exp":        time.Now().Add(time.Second * time.Duration(jwtExpiry)).Unix(),
		}).SignedString(jwtKey)
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
}

func NewGetMeHandler(config *internal.Config, users *database.Users) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("auth_user")
		if !ok {
			SendServerError(ctx, "an error occured")
			return
		}
		authUser, ok := val.(*database.User)
		if !ok {
			SendServerError(ctx, "an error occured")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": authUser,
		})
	}
}

func NewLogoutHandler(config *internal.Config, tokens *database.BlacklistedTokens) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		err := tokens.Add(authToken.Raw, int64(claims["exp"].(float64)))
		if err != nil {
			SendServerError(ctx, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}
