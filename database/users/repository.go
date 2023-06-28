package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type usersRepository struct {
	conn *pgx.Conn
}

func (users *usersRepository) Create(ctx context.Context, user *User) error {
	id := user.ID.String()
	if strings.EqualFold(id, "") {
		return fmt.Errorf("invalid user id")
	}

	res, err := users.conn.Exec(
		ctx,
		`INSERT INTO service.users (id, name, email, password) VALUES($1, $2, $3, $4);`,
		id, user.Name, user.Email, user.Password,
	)
	if err != nil {
		return err
	} else if res.RowsAffected() != 1 {
		return errors.New("unable to insert new user")
	}

	return nil
}

func (users *usersRepository) GetOneById(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{ID: id}
	err := users.conn.QueryRow(
		ctx,
		`SELECT name, email, password FROM service.users WHERE id = $1 LIMIT 1`,
		id.String(),
	).Scan(&user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (users *usersRepository) GetOneByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{Email: email}
	var id string

	err := users.conn.QueryRow(
		ctx,
		`SELECT id, name, password FROM service.users WHERE LOWER(email) = LOWER($1) LIMIT 1`,
		email,
	).Scan(&id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	user.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUsersRepository(conn *pgx.Conn) UsersRepository {
	return &usersRepository{conn}
}
