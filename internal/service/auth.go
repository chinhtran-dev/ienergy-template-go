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
	"ienergy-template-go/pkg/logger"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req request.UserLoginRequest) (resp response.TokenResponse, err error)
	Register(ctx context.Context, req request.UserRegisterRequest) (resp response.UserInfoResponse, err error)
}

type authService struct {
	userRepo repository.UserRepo
	logger   *logger.StandardLogger
	config   *config.Config
}

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

// Login implements IAuthService.
func (s *authService) Login(ctx context.Context, req request.UserLoginRequest) (resp response.TokenResponse, err error) {
	userID, err := s.userRepo.ValidateUser(entity.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.logger.WithField("err", err.Error()).Info("User")
		return
	}

	if userID == uuid.Nil {
		s.logger.Info("User")
		return
	}

	token, err := s.generateToken(userID, req.Email)
	if err != nil {
		return
	}

	resp = response.TokenResponse{
		Token: token,
	}

	return
}

// Register implements IAuthService.
func (s *authService) Register(ctx context.Context, req request.UserRegisterRequest) (resp response.UserInfoResponse, err error) {
	err = s.userRepo.VerifyUserEmail(ctx, req.Email)
	if err != nil {
		s.logger.
			WithKeyword(ctx, "VerifyUserEmail").
			WithError(err).
			Error()
		return
	}

	user, err := s.userRepo.UserRegister(ctx, entity.ToEntityModel(req))
	if err != nil {
		return
	}
	if user.ID == uuid.Nil {
		return
	}

	resp = response.UserInfoResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
	}

	return
}

func (s *authService) generateToken(userID uuid.UUID, email string) (string, error) {
	lifespan, err := strconv.Atoi(s.config.JWT.ExpirationTime)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims[constant.UserID] = userID
	claims[constant.Email] = email
	claims[constant.ExpireDate] = time.Now().Add(time.Hour * time.Duration(lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
