package rest

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
)

type RESTServer struct {
	router *gin.Engine
}

func (s *RESTServer) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}

func NewRESTServer(usersService users.UsersService) (*RESTServer, error) {
	// Create route handlers
	usersAuthHandlers := NewUsersAuthHandlers(usersService)
	middlewares := NewMiddlewares(usersService)

	// Create and configure router
	router := gin.New()
	router.Use(gin.Recovery(), cors.Default())

	// Register routes
	router.POST("/register", usersAuthHandlers.HandleRegister)
	router.POST("/generate-access-token", usersAuthHandlers.HandleGenerateAccessToken)
	router.POST("/blacklist-access-token", middlewares.HandleUserAuth, usersAuthHandlers.HandleBlacklistAccessToken)
	router.GET("/me", middlewares.HandleUserAuth, usersAuthHandlers.HandleGetMe)

	router.NoRoute(func(ctx *gin.Context) {
		SendNotFound(ctx, "The resource you were looking for was not found on this server.")
	})

	return &RESTServer{router}, nil
}
