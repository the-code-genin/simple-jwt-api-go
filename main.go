package main

import (
	"sync"

	"github.com/the-code-genin/simple-jwt-api-go/api/rest"
	"github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/postgres"
	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
	"github.com/the-code-genin/simple-jwt-api-go/database"
)

func main() {
	// Load configuration variables
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Create db connections
	pqConn, err := postgres.NewConnection(config)
	if err != nil {
		panic(err)
	}
	redisClient, err := redis.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create repositories
	usersRepo := database.NewUsersRepository(pqConn)
	blacklistedTokensRepo := database.NewBlacklistedTokensRepository(config, redisClient)

	// Create Services
	usersService := users.NewUsersService(config, usersRepo, blacklistedTokensRepo)

	// Create API servers
	server, err := rest.NewRESTServer(usersService)
	if err != nil {
		panic(err)
	}

	// Run system components
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Run(config.Port); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
