package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/the-code-genin/simple-jwt-api-go/database"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

type Server struct {
	ctx    *internal.AppContext
	router *gin.Engine
}

func (s *Server) Start() error {
	port, err := s.ctx.GetConfig().GetHTTPPort()
	if err != nil {
		return err
	}
	return s.router.Run(fmt.Sprintf(":%d", port))
}

func NewServer(ctx *internal.AppContext) *Server {
	config := ctx.GetConfig()
	users := database.NewUsers(ctx)
	blacklistedTokens := database.NewBlacklistedTokens(ctx)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/signup", NewSignupHandler(users))
	router.POST("/generate-access-token", NewLoginHandler(config, users))
	router.POST(
		"/blacklist-access-token",
		NewAuthMiddleware(config, users, blacklistedTokens),
		NewLogoutHandler(config, blacklistedTokens),
	)
	router.GET("/me", NewAuthMiddleware(config, users, blacklistedTokens), NewGetMeHandler(config, users))

	router.NoRoute(func(ctx *gin.Context) {
		SendNotFound(ctx, "The resource you were looking for was not found on this server.")
	})

	return &Server{ctx, router}
}
