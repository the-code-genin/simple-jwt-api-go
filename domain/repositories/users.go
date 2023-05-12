package repositories

import (
	"github.com/google/uuid"
	"github.com/the-code-genin/simple-jwt-api-go/domain/entities"
)

type UsersRepository interface {
	Create(user *entities.User) error

	GetOneById(id uuid.UUID) (*entities.User, error)
	GetOneByEmail(email string) (*entities.User, error)

	EmailTaken(email string) (bool, error)
}
