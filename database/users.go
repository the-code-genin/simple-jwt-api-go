package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/the-code-genin/simple-jwt-api-go/domain/entities"
	"github.com/the-code-genin/simple-jwt-api-go/domain/repositories"
)

type UsersRepository struct {
	conn *pgx.Conn
}

func (users *UsersRepository) Create(user *entities.User) error {
	id := user.ID.String()
	if strings.EqualFold(id, "") {
		return fmt.Errorf("invalid user id")
	}

	res, err := users.conn.Exec(
		context.Background(),
		`INSERT INTO users (id, name, email, password) VALUES($1, $2, LOWER($3), $4);`,
		id, user.Name, user.Email, user.Password,
	)
	if err != nil {
		return err
	} else if res.RowsAffected() != 1 {
		return errors.New("unable to insert new user")
	}

	return nil
}

func (users *UsersRepository) GetOneById(id uuid.UUID) (*entities.User, error) {
	user := &entities.User{ID: id}
	err := users.conn.QueryRow(
		context.Background(),
		`SELECT name, email, password FROM users WHERE id = $1 LIMIT 1`,
		id.String(),
	).Scan(&user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (users *UsersRepository) GetOneByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	var id string

	err := users.conn.QueryRow(
		context.Background(),
		`SELECT id, name, email, password FROM users WHERE email = LOWER($1) LIMIT 1`,
		email,
	).Scan(&id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	user.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUsersRepository(conn *pgx.Conn) repositories.UsersRepository {
	return &UsersRepository{conn}
}
