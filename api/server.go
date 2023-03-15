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
	users := database.NewUsers(ctx)

	router.POST("/signup", NewSignupHandler(users))

	return &Server{ctx, router}
}
