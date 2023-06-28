package http

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/the-code-genin/simple-jwt-api-go/services/http/docs"
	"github.com/the-code-genin/simple-jwt-api-go/services/http/handlers"
)

type Server struct {
	router *gin.Engine
}

func (s *Server) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}

// @title       Simple JWT API Go
// @version     1.0
// @description A simple JWT powered API written in Go
// @host        localhost:9000
// @BasePath    /
// @accept      json
// @produce     json
func NewServer(isProd bool, usersService users.UsersService) (*Server, error) {
	// Create route handlers
	usersFacade := handlers.NewUsersFacade(usersService)
	middlewares := handlers.NewMiddlewares(usersService)

	// Create and configure router
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery(), cors.Default())

	router.NoRoute(func(ctx *gin.Context) {
		handlers.SendNotFound(ctx, "The resource you were looking for was not found on this server.")
	})

	// Register routes
	if !isProd {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.POST("/register", usersFacade.Register)
	router.POST("/generate-access-token", usersFacade.GenerateAccessToken)
	router.POST("/blacklist-access-token", middlewares.HandleUserAuth, usersFacade.BlacklistAccessToken)
	router.GET("/me", middlewares.HandleUserAuth, usersFacade.GetMe)

	return &Server{router}, nil
}
