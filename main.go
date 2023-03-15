package main

import (
	"context"

	"github.com/the-code-genin/simple-jwt-api-go/api"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

func main() {
	ctx := internal.NewAppContext(context.Background())
	server := api.NewServer(ctx)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
