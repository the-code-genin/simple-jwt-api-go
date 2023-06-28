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

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}
