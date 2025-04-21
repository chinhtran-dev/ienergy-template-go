package service_test

import (
	"context"
	"ienergy-template-go/internal/model/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock implementation of repository.UserRepo
type MockUserRepo struct {
	mock.Mock
}

// GetUserByID implements repository.UserRepo
func (m *MockUserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (resp entity.User, err error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.User), args.Error(1)
}

// GetUserByEmail implements repository.UserRepo
func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (resp entity.User, err error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

// UserRegister implements repository.UserRepo
func (m *MockUserRepo) UserRegister(ctx context.Context, userInfo entity.User) (resp entity.User, err error) {
	args := m.Called(ctx, userInfo)
	return args.Get(0).(entity.User), args.Error(1)
}

// ValidateUser implements repository.UserRepo
func (m *MockUserRepo) ValidateUser(userInfo entity.User) (userID uuid.UUID, err error) {
	args := m.Called(userInfo)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// UpdateUser implements repository.UserRepo
func (m *MockUserRepo) UpdateUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Error(0)
}

// DeleteUser implements repository.UserRepo
func (m *MockUserRepo) DeleteUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Error(0)
}

// VerifyUserEmail implements repository.UserRepo
func (m *MockUserRepo) VerifyUserEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
