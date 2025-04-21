package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/pkg/wrapper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService implements the AuthService interface for testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req request.UserRegisterRequest) (response.UserInfoResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(response.UserInfoResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req request.UserLoginRequest) (response.TokenResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(response.TokenResponse), args.Error(1)
}

// TestAuthHandler_Register tests the Register handler functionality
func TestAuthHandler_Register(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)

	// Define test cases
	testCases := []struct {
		name         string
		req          request.UserRegisterRequest
		mockSetup    func(*MockAuthService)
		expectedCode int
		validateResp func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			req: request.UserRegisterRequest{
				Email:           "test@example.com",
				Password:        "password123",
				FirstName:       "John",
				LastName:        "Doe",
				ConfirmPassword: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.Anything).Return(
					response.UserInfoResponse{
						UserID:   uuid.New(),
						Email:    "test@example.com",
						FullName: "John Doe",
					}, nil,
				)
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NotNil(t, resp.Data)

				// Validate response data
				var userInfo response.UserInfoResponse
				dataBytes, _ := json.Marshal(resp.Data)
				err = json.Unmarshal(dataBytes, &userInfo)
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", userInfo.Email)
				assert.Equal(t, "John Doe", userInfo.FullName)
			},
		},
		{
			name: "invalid request - missing email",
			req: request.UserRegisterRequest{
				Password:        "password123",
				FirstName:       "John",
				LastName:        "Doe",
				ConfirmPassword: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed as validation should fail before service call
			},
			expectedCode: http.StatusBadRequest,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
			},
		},
		{
			name: "email already exists",
			req: request.UserRegisterRequest{
				Email:           "existing@example.com",
				Password:        "password123",
				FirstName:       "John",
				LastName:        "Doe",
				ConfirmPassword: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.Anything).Return(
					response.UserInfoResponse{},
					errors.New("email already exists"),
				)
			},
			expectedCode: http.StatusInternalServerError,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
				if resp.Message != nil {
					assert.Equal(t, "email already exists", *resp.Message)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockAuthService)

			// Create request body
			reqBody, err := json.Marshal(tc.req)
			assert.NoError(t, err)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute handler
			authHandler.Register()(c)

			// Validate response
			assert.Equal(t, tc.expectedCode, w.Code)
			tc.validateResp(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}

// TestAuthHandler_Login tests the Login handler functionality
func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)

	// Define test cases
	testCases := []struct {
		name         string
		req          request.UserLoginRequest
		mockSetup    func(*MockAuthService)
		expectedCode int
		validateResp func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			req: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.Anything).Return(
					response.TokenResponse{
						Token: "valid_token_123",
					}, nil,
				)
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NotNil(t, resp.Data)

				// Validate response data
				var tokenResp response.TokenResponse
				dataBytes, _ := json.Marshal(resp.Data)
				err = json.Unmarshal(dataBytes, &tokenResp)
				assert.NoError(t, err)
				assert.Equal(t, "valid_token_123", tokenResp.Token)
			},
		},
		{
			name: "invalid request - missing email",
			req: request.UserLoginRequest{
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed as validation should fail before service call
			},
			expectedCode: http.StatusBadRequest,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
			},
		},
		{
			name: "invalid credentials",
			req: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.Anything).Return(
					response.TokenResponse{},
					errors.New("invalid credentials"),
				)
			},
			expectedCode: http.StatusUnauthorized,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
				if resp.Message != nil {
					assert.Equal(t, "invalid credentials", *resp.Message)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockAuthService)

			// Create request body
			reqBody, err := json.Marshal(tc.req)
			assert.NoError(t, err)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute handler
			authHandler.Login()(c)

			// Validate response
			assert.Equal(t, tc.expectedCode, w.Code)
			tc.validateResp(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}
