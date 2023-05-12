package users

import "fmt"

var (
	ErrEmailTaken      = fmt.Errorf("email is taken")
	ErrInvalidPassword = fmt.Errorf("invalid password")
)
