package main

import (
	"github.com/the-code-genin/simple-jwt-api-go/api"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

func main() {
	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Start server
	server, err := api.NewServer(config)
	if err != nil {
		panic(err)
	}
	if err := server.Start(); err != nil {
		panic(err)
	}
}
