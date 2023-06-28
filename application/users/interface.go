package users

import (
	"context"
)

type UsersService interface {
	Register(ctx context.Context, req RegisterUserDTO) (*UserDTO, error)

	GenerateAccessToken(ctx context.Context, req GenerateUserAccessTokenDTO) (user *UserDTO, token string, err error)
	DecodeAccessToken(ctx context.Context, token string) (*UserDTO, error)
	BlacklistAccessToken(ctx context.Context, token string) error
}
