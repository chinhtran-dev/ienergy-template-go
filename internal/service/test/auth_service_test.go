package service_test

import (
	"context"
	"errors"
	"ienergy-template-go/config"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/constant"
	"ienergy-template-go/pkg/logger"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAuthService_Login tests the Login functionality of AuthService
func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockUserRepo := new(MockUserRepo)
	mockConfig := &config.Config{
		Server: config.ServerCfg{
			Env: constant.DevelopmentEnv,
		},
		JWT: config.JWTConfig{
			Secret:         "secret",
			ExpirationTime: "1", // 1 hour
		},
	}
	mockLogger := logger.NewLogger(mockConfig)

	authService := service.NewAuthService(mockUserRepo, mockLogger, mockConfig)

	// Define test cases
	testCases := []struct {
		name          string
		req           request.UserLoginRequest
		mockSetup     func(*MockUserRepo)
		expectedError error
		validateResp  func(*testing.T, response.TokenResponse, error)
	}{
		{
			name: "successful login",
			req: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserRepo) {
				m.On("ValidateUser", mock.Anything).Return(uuid.New(), nil)
			},
			expectedError: nil,
			validateResp: func(t *testing.T, resp response.TokenResponse, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
			},
		},
		{
			name: "invalid credentials",
			req: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserRepo) {
				m.On("ValidateUser", mock.Anything).Return(uuid.Nil, errors.New("invalid credentials"))
			},
			expectedError: errors.New("invalid credentials"),
			validateResp: func(t *testing.T, resp response.TokenResponse, err error) {
				assert.Error(t, err)
				assert.Empty(t, resp.Token)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockUserRepo)

			// Execute test
			resp, err := authService.Login(context.Background(), tc.req)

			// Validate results
			tc.validateResp(t, resp, err)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestAuthService_Register tests the Register functionality of AuthService
func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockUserRepo := new(MockUserRepo)
	mockConfig := &config.Config{
		Server: config.ServerCfg{
			Env: constant.DevelopmentEnv,
		},
		JWT: config.JWTConfig{
			Secret:         "secret",
			ExpirationTime: "1", // 1 hour
		},
	}
	mockLogger := logger.NewLogger(mockConfig)

	authService := service.NewAuthService(mockUserRepo, mockLogger, mockConfig)

	// Define test cases
	testCases := []struct {
		name          string
		req           request.UserRegisterRequest
		mockSetup     func(*MockUserRepo)
		expectedError error
		validateResp  func(*testing.T, response.UserInfoResponse, error)
	}{
		{
			name: "successful registration",
			req: request.UserRegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(m *MockUserRepo) {
				m.On("VerifyUserEmail", mock.Anything, "test@example.com").Return(nil)
				userID := uuid.New()
				user := entity.User{
					ID:        userID,
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				m.On("UserRegister", mock.Anything, mock.Anything).Return(user, nil)
			},
			expectedError: nil,
			validateResp: func(t *testing.T, resp response.UserInfoResponse, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.UserID)
				assert.Equal(t, "test@example.com", resp.Email)
				assert.Equal(t, "John Doe", resp.FullName)
			},
		},
		{
			name: "email already exists",
			req: request.UserRegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(m *MockUserRepo) {
				m.On("VerifyUserEmail", mock.Anything, "existing@example.com").Return(errors.New("email already exists"))
			},
			expectedError: errors.New("email already exists"),
			validateResp: func(t *testing.T, resp response.UserInfoResponse, err error) {
				assert.Error(t, err)
				assert.Empty(t, resp.UserID)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockUserRepo)

			// Execute test
			resp, err := authService.Register(context.Background(), tc.req)

			// Validate results
			tc.validateResp(t, resp, err)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
