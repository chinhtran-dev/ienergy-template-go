package service_test

import (
	"context"
	"errors"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/constant"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserService_GetUserInfo tests the GetUserInfo functionality of UserService
func TestUserService_GetUserInfo(t *testing.T) {
	t.Parallel()

	// Setup test dependencies
	mockUserRepo := new(MockUserRepo)
	mockDB := new(MockDatabase)
	userService := service.NewUserService(mockUserRepo, mockDB)

	// Define test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(*MockUserRepo, *MockDatabase)
		expectedError error
		validateResp  func(*testing.T, response.UserInfoResponse, error)
	}{
		{
			name:   "successful get user info",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				mockUser := entity.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				m.On("GetUserByID", mock.Anything, mock.Anything).Return(mockUser, nil)
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
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				m.On("GetUserByID", mock.Anything, mock.Anything).Return(entity.User{}, errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
			validateResp: func(t *testing.T, resp response.UserInfoResponse, err error) {
				assert.Error(t, err)
				assert.Empty(t, resp.UserID)
			},
		},
		{
			name:   "invalid user id in context",
			userID: uuid.Nil,
			mockSetup: func(m *MockUserRepo, db *MockDatabase) {
				// No mock setup needed as the function should return early
			},
			expectedError: nil,
			validateResp: func(t *testing.T, resp response.UserInfoResponse, err error) {
				assert.NoError(t, err)
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
			tc.mockSetup(mockUserRepo, mockDB)

			// Create context with userID
			ctx := context.WithValue(context.Background(), constant.UserID, tc.userID)

			// Execute test
			resp, err := userService.GetUserInfo(ctx)

			// Validate results
			tc.validateResp(t, resp, err)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
