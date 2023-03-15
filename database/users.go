package database

import (
	"context"
	"errors"

	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Users struct {
	ctx *internal.AppContext
}

// Get a single user by their ID
func (users *Users) GetOne(id int) (*User, error) {
	conn, err := users.ctx.GetPostgresConn()
	if err != nil {
		return nil, err
	}

	user := &User{}
	user.ID = id
	err = conn.QueryRow(
		context.Background(),
		`SELECT name, email, password FROM users WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create a new user
func (users *Users) Insert(user *User) (*User, error) {
	conn, err := users.ctx.GetPostgresConn()
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	// Insert user data
	res, err := tx.Exec(
		context.Background(),
		`INSERT INTO users (name, email, password) VALUES($1, $2, $3);`,
		user.Name, user.Email, user.Password,
	)
	if err != nil {
		return nil, err
	} else if res.RowsAffected() != 1 {
		return nil, errors.New("unable to insert new user")
	}

	// Get the new user id
	err = tx.QueryRow(
		context.Background(),
		`SELECT id FROM users ORDER BY id DESC LIMIT 1;`,
	).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	// Save db changes
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUsers(ctx *internal.AppContext) *Users {
	return &Users{ctx}
}
