package service_test

import (
	"context"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestUserService_GetUserInfo tests the GetUserInfo functionality of UserService
func TestUserService_GetUserInfo(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockUserRepo := new(MockUserRepo)
	mockDB := new(MockDatabase)

	// Setup mock database to return a valid DB instance
	mockDB.On("GetDB").Return(&gorm.DB{})

	userService := service.NewUserService(mockUserRepo, mockDB)

	// Define test cases
	testCases := []struct {
		name         string
		userID       uuid.UUID
		mockSetup    func(*MockUserRepo, *MockDatabase)
		expectedResp response.UserInfoResponse
		expectedErr  *errors.AppError
	}{
		{
			name:   "successful get user info",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				m.On("GetUserByID", mock.Anything, mock.Anything).Return(response.UserInfoResponse{
					UserID:   uuid.New(),
					Email:    "test@example.com",
					FullName: "John Doe",
				}, nil)
			},
			expectedResp: response.UserInfoResponse{
				UserID:   uuid.New(),
				Email:    "test@example.com",
				FullName: "John Doe",
			},
			expectedErr: nil,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				m.On("GetUserByID", mock.Anything, mock.Anything).Return(response.UserInfoResponse{},
					errors.NewNotFoundError("user not found"))
			},
			expectedResp: response.UserInfoResponse{},
			expectedErr:  errors.NewNotFoundError("user not found"),
		},
		{
			name:   "invalid user id in context",
			userID: uuid.Nil,
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				// No mock setup needed as validation should fail before service call
			},
			expectedResp: response.UserInfoResponse{},
			expectedErr:  errors.NewUnauthorizedError("invalid user id"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			tc.mockSetup(mockUserRepo, mockDB)

			// Create context with user ID
			ctx := context.Background()
			if tc.userID != uuid.Nil {
				ctx = context.WithValue(ctx, "user_id", tc.userID)
			}

			// Execute service method
			resp, err := userService.GetUserInfo(ctx)

			// Validate response
			assert.Equal(t, tc.expectedResp, resp)
			assert.Equal(t, tc.expectedErr, err)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
