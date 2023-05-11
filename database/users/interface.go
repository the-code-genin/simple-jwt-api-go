package users

import (
	"github.com/the-code-genin/simple-jwt-api-go/domain"
)

type Users interface {
	GetOneById(id int) (*domain.User, error)
	GetOneByEmail(email string) (*domain.User, error)

	Create(user *domain.User) (*domain.User, error)

	EmailTaken(email string) (bool, error)
}
