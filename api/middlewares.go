package api

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database/repositories"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

type Middlewares struct {
	config *internal.Config
	users  *repositories.Users
	tokens *repositories.BlacklistedTokens
}

func (m *Middlewares) HandleAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	match := regexp.MustCompile(`^Bearer\s+([^\s]+)$`).FindStringSubmatch(authHeader)
	if len(match) != 2 {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	}

	token, err := jwt.Parse(match[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.JWT.Key), nil
	})
	if err != nil {
		SendBadRequest(ctx, err.Error())
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	}

	userID, userIDOk := claims["user_id"].(float64)
	userEmail, userEmailOk := claims["user_email"].(string)
	exp, expOk := claims["exp"].(float64)
	if !userIDOk || !userEmailOk || !expOk {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	}

	user, err := m.users.GetOne(int(userID))
	if err != nil || user == nil {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	} else if user.Email != userEmail {
		SendBadRequest(ctx, "invalid Authorization header")
		ctx.Abort()
		return
	} else if time.Now().After(time.Unix(int64(exp), 0)) {
		SendBadRequest(ctx, "expired Authorization header")
		ctx.Abort()
		return
	}

	blacklisted, err := m.tokens.Exists(token.Raw)
	if err != nil {
		SendBadRequest(ctx, err.Error())
		ctx.Abort()
		return
	} else if blacklisted {
		SendBadRequest(ctx, "blacklisted Authorization header")
		ctx.Abort()
		return
	}

	ctx.Set("auth_user", user)
	ctx.Set("auth_token", token)
	ctx.Next()
}

func NewMiddlewares(
	config *internal.Config,
	users *repositories.Users,
	tokens *repositories.BlacklistedTokens,
) *Middlewares {
	return &Middlewares{config, users, tokens}
}
