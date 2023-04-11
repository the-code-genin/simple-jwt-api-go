package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/simple-jwt-api-go/database/repositories"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

type Server struct {
	port   int
	router *gin.Engine
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func NewServer(config *internal.Config) (*Server, error) {
	// Create db connections
	dbConn, err := internal.ConnectToPostgres(config)
	if err != nil {
		return nil, err
	}
	redisClient, err := internal.ConnectToRedis(config)
	if err != nil {
		return nil, err
	}

	// Create repositories
	users := repositories.NewUsers(dbConn)
	blacklistedTokens := repositories.NewBlacklistedTokens(config, redisClient)

	// Create handlers
	authHandlers := NewAuthHandlers(config, users, blacklistedTokens)
	middlewares := NewMiddlewares(config, users, blacklistedTokens)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/signup", authHandlers.HandleSignup)
	router.POST("/generate-access-token", authHandlers.HandleLogin)
	router.POST("/blacklist-access-token", middlewares.HandleAuth, authHandlers.HandleLogout)
	router.GET("/me", middlewares.HandleAuth, authHandlers.HandleGetMe)

	router.NoRoute(func(ctx *gin.Context) {
		SendNotFound(ctx, "The resource you were looking for was not found on this server.")
	})

	return &Server{config.Port, router}, nil
}
