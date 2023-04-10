package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/the-code-genin/simple-jwt-api-go/api"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

func main() {
	// Load env variables
	if _, err := os.Stat(".env"); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	// Parse config data
	var config internal.Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	// Start server
	server := api.NewServer(config)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
