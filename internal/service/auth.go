package service

import (
	"context"
	"fmt"
	"ienergy-template-go/config"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/internal/repository"
	"ienergy-template-go/pkg/constant"
	"ienergy-template-go/pkg/errors"
	"ienergy-template-go/pkg/logger"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, req request.UserLoginRequest) (response.TokenResponse, error)
	Register(ctx context.Context, req request.UserRegisterRequest) (response.UserInfoResponse, error)
}

// authService implements AuthService
type authService struct {
	userRepo repository.UserRepo
	logger   *logger.StandardLogger
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepo,
	logger *logger.StandardLogger,
	config *config.Config,
) AuthService {
	return &authService{
		userRepo: userRepo,
		logger:   logger,
		config:   config,
	}
}

// Login handles user login
func (s *authService) Login(ctx context.Context, req request.UserLoginRequest) (response.TokenResponse, error) {
	userID, err := s.userRepo.ValidateUser(entity.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.logger.WithField("err", err.Error()).Info("Login failed")
		return response.TokenResponse{}, err
	}

	if userID == uuid.Nil {
		s.logger.Info("Invalid credentials")
		return response.TokenResponse{}, errors.NewUnauthorizedError("Invalid email or password")
	}

	token, tokenErr := s.generateToken(userID, req.Email)
	if tokenErr != nil {
		s.logger.WithError(tokenErr).Error("Failed to generate token")
		return response.TokenResponse{}, errors.NewInternalServerError("Failed to generate token: " + tokenErr.Error())
	}

	return response.TokenResponse{
		Token: token,
	}, nil
}

// Register handles user registration
func (s *authService) Register(ctx context.Context, req request.UserRegisterRequest) (response.UserInfoResponse, error) {
	err := s.userRepo.VerifyUserEmail(ctx, req.Email)
	if err != nil {
		s.logger.
			WithContext(ctx).
			WithField("email", req.Email).
			WithError(err).
			Error("Email verification failed")
		return response.UserInfoResponse{}, errors.NewConflictError("email already exists")
	}

	user, err := s.userRepo.UserRegister(ctx, entity.ToEntityModel(req))
	if err != nil {
		s.logger.
			WithContext(ctx).
			WithField("email", req.Email).
			WithError(err).
			Error("User registration failed")
		return response.UserInfoResponse{}, errors.NewInternalServerError("failed to create user")
	}

	if user.ID == uuid.Nil {
		s.logger.
			WithContext(ctx).
			WithField("email", req.Email).
			Error("User registration failed - invalid user ID")
		return response.UserInfoResponse{}, errors.NewBadRequestError("Failed to create user")
	}

	return response.UserInfoResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
	}, nil
}

// generateToken generates JWT token
func (s *authService) generateToken(userID uuid.UUID, email string) (string, *errors.AppError) {
	lifespan, err := strconv.Atoi(s.config.JWT.ExpirationTime)
	if err != nil {
		return "", errors.NewUnauthorizedError("Invalid token expiration time: " + err.Error())
	}

	claims := jwt.MapClaims{
		constant.UserID:     userID,
		constant.Email:      email,
		constant.ExpireDate: time.Now().Add(time.Hour * time.Duration(lifespan)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", errors.NewUnauthorizedError("Failed to sign token: " + err.Error())
	}

	return tokenString, nil
}
