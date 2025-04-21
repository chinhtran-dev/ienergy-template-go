package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAuthHandler_Register(t *testing.T) {
	// Mock AuthService
	mockAuthService := new(MockAuthService)

	// Tạo AuthHandler với mock AuthService
	handler := handler.NewAuthHandler(mockAuthService)

	// Định nghĩa request giả
	registerReq := request.UserRegisterRequest{
		Email:           "test@example.com",
		Password:        "password123",
		FirstName:       "John",
		LastName:        "Doe",
		ConfirmPassword: "password123",
	}

	// Mock hàm Register trong AuthService
	mockAuthService.On("Register", mock.Anything, registerReq).Return(
		response.UserInfoResponse{
			UserID:   uuid.New(),
			Email:    "test@example.com",
			FullName: "John Doe",
		}, nil,
	).Once()

	// Tạo request JSON cho đăng ký
	reqBody, err := json.Marshal(registerReq)
	assert.NoError(t, err)

	// Tạo HTTP request và recorder để ghi nhận kết quả
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Tạo gin context và gọi handler
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req
	handler.Register()(c)

	// Kiểm tra status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Kiểm tra response body
	var resp wrapper.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthHandler_Register_InvalidRequest(t *testing.T) {
	// Mock AuthService
	mockAuthService := new(MockAuthService)

	// Tạo AuthHandler với mock AuthService
	handler := handler.NewAuthHandler(mockAuthService)

	// Tạo request giả không hợp lệ (thiếu email)
	registerReq := request.UserRegisterRequest{
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Tạo request JSON cho đăng ký
	reqBody, err := json.Marshal(registerReq)
	assert.NoError(t, err)

	// Tạo HTTP request và recorder để ghi nhận kết quả
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Tạo gin context và gọi handler
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req
	handler.Register()(c)

	// Kiểm tra status code lỗi
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Kiểm tra response body
	var resp wrapper.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Message)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuthHandler_Login(t *testing.T) {
	// Mock AuthService
	mockAuthService := new(MockAuthService)

	// Tạo AuthHandler với mock AuthService
	handler := handler.NewAuthHandler(mockAuthService)

	// Tạo request giả cho đăng nhập
	loginReq := request.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock hàm Login trong AuthService
	mockAuthService.On("Login", mock.Anything, loginReq).Return(
		response.TokenResponse{
			Token: "valid_token_123",
		}, nil,
	).Once()

	// Tạo request JSON cho đăng nhập
	reqBody, err := json.Marshal(loginReq)
	assert.NoError(t, err)

	// Tạo HTTP request và recorder để ghi nhận kết quả
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Tạo gin context và gọi handler
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req
	handler.Login()(c)

	// Kiểm tra status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Kiểm tra response body
	var resp wrapper.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	// Mock AuthService
	mockAuthService := new(MockAuthService)

	// Tạo AuthHandler với mock AuthService
	handler := handler.NewAuthHandler(mockAuthService)

	// Tạo request giả không hợp lệ (thiếu email)
	loginReq := request.UserLoginRequest{
		Password: "password123",
	}

	// Tạo request JSON cho đăng nhập
	reqBody, err := json.Marshal(loginReq)
	assert.NoError(t, err)

	// Tạo HTTP request và recorder để ghi nhận kết quả
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Tạo gin context và gọi handler
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req
	handler.Login()(c)

	// Kiểm tra status code lỗi
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Kiểm tra response body
	var resp wrapper.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Message)
	assert.Equal(t, 400, resp.StatusCode)
}
