package users

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"github.com/the-code-genin/simple-jwt-api-go/database/blacklisted_tokens"
	"github.com/the-code-genin/simple-jwt-api-go/database/users"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	config                      *config.Config
	usersRepository             users.UsersRepository
	blacklistedTokensRepository blacklisted_tokens.BlacklistedTokensRepository
}

func (s *usersService) Register(ctx context.Context, req RegisterUserDTO) (*UserDTO, error) {
	ctx = logger.With(ctx, zap.String(logger.FunctionNameField, "UsersService/Register"))

	// Check if the email is taken
	existingUser, err := s.usersRepository.GetOneByEmail(ctx, req.Email)
	if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		logger.Error(ctx, "An error occured while getting the user by email", zap.Error(err))
		return nil, err
	}

	if existingUser != nil {
		err := errors.New("email taken")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		logger.Error(ctx, "An error occured while hashing user password", zap.Error(err))
		return nil, err
	}

	// Create the user record
	user := &users.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hex.EncodeToString(hashedPassword),
	}
	if err := s.usersRepository.Create(ctx, user); err != nil {
		logger.Error(ctx, "An error occured while creating user", zap.Error(err))
		return nil, err
	}

	return parseUserToUserDTO(user)
}

func (s *usersService) GenerateAccessToken(ctx context.Context, req GenerateUserAccessTokenDTO) (*UserDTO, string, error) {
	ctx = logger.With(ctx, zap.String(logger.FunctionNameField, "UsersService/GenerateAccessToken"))

	// Get the user and verify the password
	user, err := s.usersRepository.GetOneByEmail(ctx, req.Email)
	if err != nil {
		logger.Error(ctx, "An error occured while getting the user by email", zap.Error(err))
		return nil, "", err
	}

	hashedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		logger.Error(ctx, "An error occured while decoding the user password", zap.Error(err))
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password)); err != nil {
		logger.Error(ctx, "Error while comparing hash and password", zap.Error(err))
		return nil, "", err
	}

	// Generate JWT token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID.String(),
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Second * time.Duration(s.config.JWT.Exp)).Unix(),
	}).SignedString([]byte(s.config.JWT.Key))
	if err != nil {
		logger.Error(ctx, "Unable to generate token for user", zap.Error(err))
		return nil, "", err
	}

	dto, err := parseUserToUserDTO(user)
	if err != nil {
		logger.Error(ctx, "Unable to parse user DTO", zap.Error(err))
		return nil, "", err
	}

	return dto, token, nil
}

func (s *usersService) DecodeAccessToken(ctx context.Context, token string) (*UserDTO, error) {
	ctx = logger.With(ctx, zap.String(logger.FunctionNameField, "UsersService/DecodeAccessToken"))

	// Parse JWT token
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Key), nil
	})
	if err != nil {
		logger.Error(ctx, "An error occured while parsing the JWT token", zap.Error(err))
		return nil, err
	}

	// Verify JWT token claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		err := errors.New("invalid JWT claims")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	userID, userIDOk := claims["user_id"].(string)
	userEmail, userEmailOk := claims["user_email"].(string)
	exp, expOk := claims["exp"].(float64)
	if !userIDOk || !userEmailOk || !expOk {
		err := errors.New("invalid/incomplete JWT claims")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error(ctx, "An error occured while parsing the userID from JWT token", zap.Error(err))
		return nil, err
	}

	// Get and verify user encoded in JWT token
	user, err := s.usersRepository.GetOneById(ctx, userUUID)
	if err != nil {
		logger.Error(ctx, "An error occured while getting the user by UUID", zap.Error(err))
		return nil, err
	}

	if !strings.EqualFold(user.Email, userEmail) {
		err := errors.New("invalid JWT claims, email doesn't match")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	if time.Now().After(time.Unix(int64(exp), 0)) {
		err := errors.New("expired access token")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	// Ensure token is not blacklisted
	blacklisted, err := s.blacklistedTokensRepository.Exists(ctx, token)
	if err != nil {
		logger.Error(ctx, "An error occured while checking blacklisted token existence", zap.Error(err))
		return nil, err
	}

	if blacklisted {
		err := errors.New("blacklisted access token")
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return parseUserToUserDTO(user)
}

func (s *usersService) BlacklistAccessToken(ctx context.Context, token string) error {
	ctx = logger.With(ctx, zap.String(logger.FunctionNameField, "UsersService/BlacklistAccessToken"))

	err := s.blacklistedTokensRepository.Add(ctx, token, int64(s.config.JWT.Exp))
	if err != nil {
		logger.Error(ctx, "Unable to blacklist access token", zap.Error(err))
		return err
	}

	return nil
}

func NewUsersService(
	config *config.Config,
	usersRepository users.UsersRepository,
	blacklistedTokensRepository blacklisted_tokens.BlacklistedTokensRepository,
) UsersService {
	return &usersService{config, usersRepository, blacklistedTokensRepository}
}
