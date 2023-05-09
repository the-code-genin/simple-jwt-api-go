package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/the-code-genin/simple-jwt-api-go/database/blacklisted_tokens"
	"github.com/the-code-genin/simple-jwt-api-go/database/users"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	config *internal.Config
	users  users.Users
	tokens blacklisted_tokens.BlacklistedTokens
}

func (s *UsersService) Register(ctx context.Context, req RegisterUserDTO) (*users.User, error) {
	// Check if the email is taken
	emailTaken, err := s.users.EmailTaken(req.Email)
	if err != nil {
		return nil, err
	}
	if emailTaken {
		return nil, ErrEmailTaken
	}

	// Hash the user's password
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	// Create the user record
	user, err := s.users.Create(&users.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hex.EncodeToString(password),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UsersService) GenerateAccessToken(ctx context.Context, req GenerateUserAccessTokenDTO) (*users.User, string, error) {
	// Get the user and verify the password
	user, err := s.users.GetOneByEmail(req.Email)
	if err != nil {
		return nil, "", err
	}
	hashedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		return nil, "", err
	}
	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	// Generate JWT token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Second * time.Duration(s.config.JWT.Exp)).Unix(),
	}).SignedString([]byte(s.config.JWT.Key))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UsersService) DecodeAccessToken(ctx context.Context, payload string) (*users.User, error) {
	// Parse access token
	token, err := jwt.Parse(payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Key), nil
	})
	if err != nil {
		return nil, err
	}

	// Verify access token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT claims")
	}
	userID, userIDOk := claims["user_id"].(float64)
	userEmail, userEmailOk := claims["user_email"].(string)
	exp, expOk := claims["exp"].(float64)
	if !userIDOk || !userEmailOk || !expOk {
		return nil, fmt.Errorf("invalid/incomplete JWT claims")
	}

	// Verify user encoded in access token
	user, err := s.users.GetOneById(int(userID))
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if user.Email != userEmail {
		return nil, fmt.Errorf("invalid auth user")
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return nil, fmt.Errorf("expired access token")
	}

	// Ensure token is not blacklisted
	blacklisted, err := s.tokens.Exists(token.Raw)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, fmt.Errorf("blacklisted access token")
	}

	return user, nil
}

func (s *UsersService) BlacklistAccessToken(ctx context.Context, payload string) error {
	err := s.tokens.Add(payload, int64(s.config.JWT.Exp))
	if err != nil {
		return err
	}

	return nil
}

func NewUsersService(
	config *internal.Config,
	users users.Users,
	tokens blacklisted_tokens.BlacklistedTokens,
) *UsersService {
	return &UsersService{config, users, tokens}
}
