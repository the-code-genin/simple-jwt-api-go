package services

import (
	"context"

	"github.com/the-code-genin/simple-jwt-api-go/domain/dtos"
	"github.com/the-code-genin/simple-jwt-api-go/domain/entities"
)

type UsersService interface {
	Register(ctx context.Context, req dtos.RegisterUserDTO) (*entities.User, error)

	GenerateAccessToken(ctx context.Context, req dtos.GenerateUserAccessTokenDTO) (*entities.User, string, error)
	DecodeAccessToken(ctx context.Context, payload string) (*entities.User, error)
	BlacklistAccessToken(ctx context.Context, payload string) error
}
