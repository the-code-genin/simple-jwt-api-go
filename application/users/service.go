package users

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/errors"
	"github.com/the-code-genin/simple-jwt-api-go/common/logger"
	"github.com/the-code-genin/simple-jwt-api-go/domain/entities"
	"github.com/the-code-genin/simple-jwt-api-go/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	config                      *config.Config
	usersRepository             repositories.UsersRepository
	blacklistedTokensRepository repositories.BlacklistedTokensRepository
}

func (s *usersService) Register(ctx context.Context, req RegisterUserDTO) (UserDTO, error) {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "usersService/Register").
		WithField(logger.RequestBodyField, req)

	// Check if the email is taken
	existingUser, err := s.usersRepository.GetOneByEmail(req.Email)
	if err != nil && !errors.IsNoRecordError(err) {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	} else if existingUser != nil {
		log.WithError(ErrEmailTaken).Error(ErrEmailTaken.Error())
		return UserDTO{}, ErrEmailTaken
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	}

	// Create the user record
	user := &entities.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hex.EncodeToString(hashedPassword),
	}
	if err := s.usersRepository.Create(user); err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	}

	return parseUserEntityToUserDTO(user)
}

func (s *usersService) GenerateAccessToken(ctx context.Context, req GenerateUserAccessTokenDTO) (UserDTO, string, error) {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "usersService/GenerateAccessToken").
		WithField(logger.RequestBodyField, req)

	// Get the user and verify the password
	user, err := s.usersRepository.GetOneByEmail(req.Email)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, "", err
	}
	hashedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, "", err
	}
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password)); err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, "", ErrInvalidPassword
	}

	// Generate JWT token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID.String(),
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Second * time.Duration(s.config.JWT.Exp)).Unix(),
	}).SignedString([]byte(s.config.JWT.Key))
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, "", err
	}

	dto, err := parseUserEntityToUserDTO(user)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, "", err
	}

	return dto, token, nil
}

func (s *usersService) DecodeAccessToken(ctx context.Context, token string) (UserDTO, error) {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "usersService/DecodeAccessToken").
		WithField(logger.TokenField, token)

	// Parse JWT token
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Key), nil
	})
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	}

	// Verify JWT token claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		err := fmt.Errorf("invalid JWT claims")
		log.Error(err.Error())
		return UserDTO{}, err
	}

	userID, userIDOk := claims["user_id"].(string)
	userEmail, userEmailOk := claims["user_email"].(string)
	exp, expOk := claims["exp"].(float64)
	if !userIDOk || !userEmailOk || !expOk {
		err := fmt.Errorf("invalid/incomplete JWT claims")
		log.Error(err.Error())
		return UserDTO{}, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	}

	// Get and verify user encoded in JWT token
	user, err := s.usersRepository.GetOneById(userUUID)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, fmt.Errorf("user not found")
	}
	if user.Email != userEmail {
		err := fmt.Errorf("invalid JWT claims")
		log.Error(err.Error())
		return UserDTO{}, err
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		err := fmt.Errorf("expired access token")
		log.Error(err.Error())
		return UserDTO{}, err
	}

	// Ensure token is not blacklisted
	blacklisted, err := s.blacklistedTokensRepository.Exists(token)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return UserDTO{}, err
	}
	if blacklisted {
		err := fmt.Errorf("blacklisted access token")
		log.Error(err.Error())
		return UserDTO{}, err
	}

	return parseUserEntityToUserDTO(user)
}

func (s *usersService) BlacklistAccessToken(ctx context.Context, token string) error {
	log := logger.NewLogger(ctx).
		WithField(logger.FunctionNameField, "usersService/BlacklistAccessToken").
		WithField(logger.TokenField, token)

	err := s.blacklistedTokensRepository.Add(token, int64(s.config.JWT.Exp))
	if err != nil {
		log.WithError(err).Error(err.Error())
		return err
	}

	return nil
}

func NewUsersService(
	config *config.Config,
	usersRepository repositories.UsersRepository,
	blacklistedTokensRepository repositories.BlacklistedTokensRepository,
) UsersService {
	return &usersService{config, usersRepository, blacklistedTokensRepository}
}
