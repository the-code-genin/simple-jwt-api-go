package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersRepository interface {
	Create(ctx context.Context, user *User) error

	GetOneById(ctx context.Context, id uuid.UUID) (*User, error)
	GetOneByEmail(ctx context.Context, email string) (*User, error)
}
