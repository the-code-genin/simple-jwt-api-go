package users

import (
	"github.com/google/uuid"
)

type UsersRepository interface {
	Create(user *User) error

	GetOneById(id uuid.UUID) (*User, error)
	GetOneByEmail(email string) (*User, error)
}
