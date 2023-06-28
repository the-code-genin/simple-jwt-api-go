package main

import (
	"context"
	"sync"

	app_users "github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"github.com/the-code-genin/simple-jwt-api-go/common/postgres"
	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
	"github.com/the-code-genin/simple-jwt-api-go/database/blacklisted_tokens"
	db_users "github.com/the-code-genin/simple-jwt-api-go/database/users"
	"github.com/the-code-genin/simple-jwt-api-go/services/http"
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

	// Create system services
	httpServer, err := http.NewServer(config.IsProduction(), usersService)
	if err != nil {
		panic(err)
	}
	log.Info("Created HTTP server")

	// Run system services
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Running HTTP server")
		if err := httpServer.Run(config.Port); err != nil {
			panic(err)
		}
	}()

	log.Info("System setup completed")
	wg.Wait()
}
