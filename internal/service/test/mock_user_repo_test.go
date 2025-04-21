package service_test

import (
	"context"
	"ienergy-template-go/internal/model/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.User), args.Get(1).(error)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Get(1).(error)
}

func (m *MockUserRepo) UserRegister(ctx context.Context, userInfo entity.User) (entity.User, error) {
	args := m.Called(ctx, userInfo)
	return args.Get(0).(entity.User), args.Get(1).(error)
}

func (m *MockUserRepo) ValidateUser(userInfo entity.User) (uuid.UUID, error) {
	args := m.Called(userInfo)
	return args.Get(0).(uuid.UUID), args.Get(1).(error)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Get(0).(error)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userInfo entity.User) error {
	args := m.Called(ctx, userInfo)
	return args.Get(0).(error)
}

func (m *MockUserRepo) VerifyUserEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Get(0).(error)
}
