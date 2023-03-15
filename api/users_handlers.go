package api

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/database"
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

		// Hash the user's password
		password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			SendBadRequest(ctx, err.Error())
			return
		}

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
