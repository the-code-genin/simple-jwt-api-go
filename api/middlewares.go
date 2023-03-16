package api

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

func NewAuthMiddleware(config *internal.Config, users *database.Users) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		re, err := regexp.Compile(`^Bearer\s+([^\s]+)$`)
		if err != nil {
			SendServerError(ctx, err.Error())
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		match := re.FindStringSubmatch(authHeader)
		if len(match) != 2 {
			SendServerError(ctx, "invalid Authorization header")
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(match[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.GetJWTKey()
		})
		if err != nil {
			SendServerError(ctx, err.Error())
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			SendServerError(ctx, "invalid Authorization header")
			ctx.Abort()
			return
		}

		userID, userIDOk := claims["user_id"].(float64)
		userEmail, userEmailOk := claims["user_email"].(string)
		exp, expOk := claims["exp"].(float64)
		if !userIDOk || !userEmailOk || !expOk {
			SendServerError(ctx, "invalid Authorization header")
			ctx.Abort()
			return
		}

		user, err := users.GetOne(int(userID))
		if err != nil || user == nil {
			SendServerError(ctx, "invalid Authorization header")
			ctx.Abort()
			return
		} else if user.Email != userEmail {
			SendServerError(ctx, "invalid Authorization header")
			ctx.Abort()
			return
		} else if time.Now().After(time.Unix(int64(exp), 0)) {
			SendServerError(ctx, "expired Authorization header")
			ctx.Abort()
			return
		}

		ctx.Set("auth_user", user)
		ctx.Next()
	}
}
