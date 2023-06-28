package users

import (
	"fmt"
	"strings"

	"github.com/the-code-genin/simple-jwt-api-go/database/users"
)

type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type GenerateUserAccessTokenDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func parseUserToUserDTO(entity *users.User) (*UserDTO, error) {
	dto := UserDTO{
		ID:    entity.ID.String(),
		Name:  entity.Name,
		Email: entity.Email,
	}

	if strings.EqualFold(dto.ID, "") {
		return nil, fmt.Errorf("invalid user id")
	}

	return &dto, nil
}
