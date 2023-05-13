package main

import (
	"context"
	"sync"

	"github.com/the-code-genin/simple-jwt-api-go/api/rest"
	app_users "github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"github.com/the-code-genin/simple-jwt-api-go/common/postgres"
	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
	"github.com/the-code-genin/simple-jwt-api-go/database/blacklisted_tokens"
	db_users "github.com/the-code-genin/simple-jwt-api-go/database/users"
)

func main() {
	log := logger.NewLogger(context.Background())

	// Load configuration variables
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Info("Loaded env variables")

	// Create db connections
	pqConn, err := postgres.NewConnection(config)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to postgres")

	redisClient, err := redis.NewClient(config)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to redis")

	// Create db repositories
	usersRepo := db_users.NewUsersRepository(pqConn)
	blacklistedTokensRepo := blacklisted_tokens.NewBlacklistedTokensRepository(config, redisClient)

	// Create application services
	usersService := app_users.NewUsersService(config, usersRepo, blacklistedTokensRepo)

	// Create API servers
	server, err := rest.NewRESTServer(config.Env, usersService)
	if err != nil {
		panic(err)
	}
	log.Info("Created REST server")

	// Run system components
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Running REST server")
		if err := server.Run(config.Port); err != nil {
			panic(err)
		}
	}()

	log.Info("System setup completed")
	wg.Wait()
}
