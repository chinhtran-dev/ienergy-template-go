package service_test

import (
	"context"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/constant"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) UserRegister(ctx context.Context, userInfo entity.User) (entity.User, error) {
	args := m.Called(ctx, userInfo)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) ValidateUser(userInfo entity.User) (uuid.UUID, error) {
	args := m.Called(userInfo)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Error(0)
}

func (m *MockUserRepo) VerifyUserEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type MockDatabase struct {
	mock.Mock
}

// GetDB implements database.Database.
func (m *MockDatabase) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// RollbackTransaction implements database.Database.
func (m *MockDatabase) RollbackTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockDatabase) BeginTransaction() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockDatabase) CommitTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockDatabase) ReleaseTransaction(tx *gorm.DB, err error) {
	m.Called(tx, err)
}

func TestGetUserInfo(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockDB := new(MockDatabase)

	// UserService với mock dependencies
	userService := service.NewUserService(mockUserRepo, mockDB)

	// Tạo một user giả để trả về
	userID := uuid.New()
	mockUser := entity.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock GetUserByID trả về thông tin người dùng
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(mockUser, nil)

	// Tạo context với userID
	ctx := context.WithValue(context.Background(), constant.UserID, userID)

	// Gọi phương thức GetUserInfo
	userInfo, err := userService.GetUserInfo(ctx)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, userID, userInfo.UserID)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "John Doe", userInfo.FullName)

	// Kiểm tra gọi phương thức GetUserByID một lần
	mockUserRepo.AssertExpectations(t)
}
