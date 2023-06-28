package main

import (
	"context"
	"os"
	"sync"

	app_users "github.com/the-code-genin/simple-jwt-api-go/application/users"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"github.com/the-code-genin/simple-jwt-api-go/common/postgres"
	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
	"github.com/the-code-genin/simple-jwt-api-go/database/blacklisted_tokens"
	db_users "github.com/the-code-genin/simple-jwt-api-go/database/users"
	"github.com/the-code-genin/simple-jwt-api-go/services/http"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Load configuration variables
	config, err := config.LoadConfig()
	if err != nil {
		logger.Error(ctx, "An error occured while loading config", zap.Error(err))
		os.Exit(1)
	}
	logger.Info(ctx, "Loaded env variables")

	// Create db connections
	pqConn, err := postgres.NewConnection(config.DB)
	if err != nil {
		logger.Error(ctx, "An error occured while connecting to the postgres database", zap.Error(err))
		os.Exit(1)
	}
	logger.Info(ctx, "Connected to postgres database")

	redisClient, err := redis.NewClient(context.Background(), config.Redis)
	if err != nil {
		logger.Error(ctx, "An error occured while connecting to redis", zap.Error(err))
		os.Exit(1)
	}
	logger.Info(ctx, "Connected to redis")

	// Create db repositories
	usersRepo := db_users.NewUsersRepository(pqConn)
	blacklistedTokensRepo := blacklisted_tokens.NewBlacklistedTokensRepository(redisClient)

	// Create application services
	usersService := app_users.NewUsersService(config, usersRepo, blacklistedTokensRepo)

	// Create system services
	httpServer, err := http.NewServer(config.IsProduction(), usersService)
	if err != nil {
		logger.Error(ctx, "An error occured while creating http server", zap.Error(err))
		os.Exit(1)
	}
	logger.Info(ctx, "Created HTTP server")

	// Run system services
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Run(config.Port); err != nil {
			logger.Error(ctx, "An error occured while running http server", zap.Error(err))
		}
	}()

	wg.Wait()
}
