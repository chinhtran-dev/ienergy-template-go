package handler

import (
	"context"
	"encoding/json"
	"ienergy-template-go/internal/http/handler"
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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserInfo(ctx context.Context) (response.UserInfoResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return response.UserInfoResponse{}, args.Error(1)
	}
	return args.Get(0).(response.UserInfoResponse), args.Error(1)
}

func TestUserHandler_Info(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	// Create test UUID
	testUUID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success - Get user info",
			setupMock: func() {
				mockService.On("GetUserInfo", mock.Anything).Return(response.UserInfoResponse{
					UserID:   testUUID,
					Email:    "test@example.com",
					FullName: "Test User",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: wrapper.NewResponse(http.StatusOK, 0, response.UserInfoResponse{
				UserID:   testUUID,
				Email:    "test@example.com",
				FullName: "Test User",
			}, "Success"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.setupMock()

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Call handler
			userHandler.Info()(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Convert expected body to JSON string
			expectedJSON, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			// Compare JSON strings
			assert.JSONEq(t, string(expectedJSON), w.Body.String())

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
