package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	config := ctx.GetConfig()
	users := database.NewUsers(ctx)

	router.POST("/signup", NewSignupHandler(users))
	router.POST("/generate-access-token", NewLoginHandler(config, users))
	router.GET("/me", NewAuthMiddleware(config, users), NewGetMeHandler(config, users))

	return &Server{ctx, router}
}
