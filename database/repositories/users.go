package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database/entities"
)

type Users struct {
	conn *pgx.Conn
}

// Get a single user by their ID
func (users *Users) GetOne(id int) (*entities.User, error) {
	user := &entities.User{}
	user.ID = id
	err := users.conn.QueryRow(
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
func (users *Users) Insert(user *entities.User) (*entities.User, error) {
	// Start transaction
	tx, err := users.conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	// Insert user data
	res, err := tx.Exec(
		context.Background(),
		`INSERT INTO users (name, email, password) VALUES($1, LOWER($2), $3);`,
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
		`SELECT id, email FROM users ORDER BY id DESC LIMIT 1;`,
	).Scan(&user.ID, &user.Email)
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

// Check if the email is taken
func (users *Users) EmailTaken(email string) (bool, error) {
	// Check if at least one user has the email
	var count int
	err := users.conn.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM users WHERE email = LOWER($1) LIMIT 1;`,
		email,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// Get the user with the email
func (users *Users) GetUserWithEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	err := users.conn.QueryRow(
		context.Background(),
		`SELECT id, name, email, password FROM users WHERE email = LOWER($1) LIMIT 1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUsers(conn *pgx.Conn) *Users {
	return &Users{conn}
}
