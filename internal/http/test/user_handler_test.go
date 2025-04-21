package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/pkg/constant"
	"ienergy-template-go/pkg/wrapper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService implements the UserService interface for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserInfo(ctx context.Context) (response.UserInfoResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(response.UserInfoResponse), args.Error(1)
}

// TestUserHandler_Info tests the Info handler functionality
func TestUserHandler_Info(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockUserService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockUserService)

	// Define test cases
	testCases := []struct {
		name         string
		userID       uuid.UUID
		mockSetup    func(*MockUserService)
		expectedCode int
		validateResp func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful get user info",
			userID: uuid.New(),
			mockSetup: func(m *MockUserService) {
				m.On("GetUserInfo", mock.Anything).Return(response.UserInfoResponse{
					UserID:   uuid.New(),
					Email:    "test@example.com",
					FullName: "John Doe",
				}, nil)
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
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *MockUserService) {
				m.On("GetUserInfo", mock.Anything).Return(response.UserInfoResponse{},
					errors.New("user not found"))
			},
			expectedCode: http.StatusNotFound,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
				if resp.Message != nil {
					assert.Equal(t, "user not found", *resp.Message)
				}
			},
		},
		{
			name:   "invalid user id in context",
			userID: uuid.Nil,
			mockSetup: func(m *MockUserService) {
				// Even though we expect an early return, we still need to set up the mock
				// to avoid unexpected method call errors
				m.On("GetUserInfo", mock.Anything).Return(response.UserInfoResponse{},
					errors.New("user not found"))
			},
			expectedCode: http.StatusNotFound,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockUserService)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Add user ID to context if not nil
			if tc.userID != uuid.Nil {
				c.Request = httptest.NewRequest("GET", "/user/info", nil)
				c.Set(constant.UserID, tc.userID)
			}

			// Execute handler
			userHandler.Info()(c)

			// Validate response
			assert.Equal(t, tc.expectedCode, w.Code)
			tc.validateResp(t, w)
			mockUserService.AssertExpectations(t)
		})
	}
}
