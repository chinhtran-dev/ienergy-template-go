package service_test

import (
	"context"
	"ienergy-template-go/config"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/logger"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockLogger := new(logger.StandardLogger)
	mockConfig := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "secret",
			ExpirationTime: "1", // 1 hour
		},
	}

	authService := service.NewAuthService(mockUserRepo, mockLogger, mockConfig)

	// Giả lập một user hợp lệ
	userID := uuid.New()
	req := request.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock ValidateUser trả về userID hợp lệ
	mockUserRepo.On("ValidateUser", mock.Anything).Return(userID, nil)

	// Gọi phương thức Login
	resp, err := authService.Login(context.Background(), req)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token) // Kiểm tra token được tạo
	mockUserRepo.AssertExpectations(t)
}

func TestRegister(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockLogger := new(logger.StandardLogger)
	mockConfig := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "secret",
			ExpirationTime: "1", // 1 hour
		},
	}

	authService := service.NewAuthService(mockUserRepo, mockLogger, mockConfig)

	// Giả lập một yêu cầu đăng ký
	req := request.UserRegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock VerifyUserEmail không trả về lỗi
	mockUserRepo.On("VerifyUserEmail", mock.Anything, req.Email).Return(nil)

	// Giả lập người dùng được đăng ký thành công
	userID := uuid.New()
	user := entity.User{
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	mockUserRepo.On("UserRegister", mock.Anything, mock.Anything).Return(user, nil)

	// Gọi phương thức Register
	resp, err := authService.Register(context.Background(), req)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, userID, resp.UserID)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "John Doe", resp.FullName)
	mockUserRepo.AssertExpectations(t)
}
