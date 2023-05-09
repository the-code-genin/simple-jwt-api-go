package main

import (
	"github.com/the-code-genin/simple-jwt-api-go/api"
	"github.com/the-code-genin/simple-jwt-api-go/database/repositories"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
	"github.com/the-code-genin/simple-jwt-api-go/services"
)

func main() {
	// Load configuration variables
	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Create db connections
	dbConn, err := internal.ConnectToPostgres(config)
	if err != nil {
		panic(err)
	}
	redisClient, err := internal.ConnectToRedis(config)
	if err != nil {
		panic(err)
	}

	// Create repositories
	users := repositories.NewUsers(dbConn)
	blacklistedTokens := repositories.NewBlacklistedTokens(config, redisClient)

	// Create Services
	usersService := services.NewUsersService(config, users, blacklistedTokens)

	// Start server
	server, err := api.NewServer(usersService)
	if err != nil {
		panic(err)
	}
	if err := server.Run(config.Port); err != nil {
		panic(err)
	}
}
