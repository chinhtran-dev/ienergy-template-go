package service_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabase is a mock implementation of database.Database
type MockDatabase struct {
	mock.Mock
}

// GetDB implements database.Database
func (m *MockDatabase) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// RollbackTransaction implements database.Database
func (m *MockDatabase) RollbackTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

// BeginTransaction implements database.Database
func (m *MockDatabase) BeginTransaction() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// CommitTransaction implements database.Database
func (m *MockDatabase) CommitTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

// ReleaseTransaction implements database.Database
func (m *MockDatabase) ReleaseTransaction(tx *gorm.DB, err error) {
	m.Called(tx, err)
}
