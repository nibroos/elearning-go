package mocks

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of the *gorm.DB object
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Commit() *gorm.DB {
	m.Called()
	return &gorm.DB{}
}

func (m *MockDB) Rollback() *gorm.DB {
	m.Called()
	return &gorm.DB{}
}

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context, filters map[string]string) ([]dtos.UserListDTO, int, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]dtos.UserListDTO), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, params *dtos.GetUserByIDParams) (*dtos.UserDetailDTO, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dtos.UserDetailDTO), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*dtos.UserDetailDTO, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*dtos.UserDetailDTO), args.Error(1)
}

func (m *MockUserRepository) BeginTransaction() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockUserRepository) AttachRoles(tx *gorm.DB, user *models.User, roleIDs []uint32) error {
	args := m.Called(tx, user, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(tx *gorm.DB, user *models.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(tx *gorm.DB, id uint) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteRolesByUserID(tx *gorm.DB, userID uint) error {
	args := m.Called(tx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) RestoreUser(tx *gorm.DB, id uint) error {
	args := m.Called(tx, id)
	return args.Error(0)
}
