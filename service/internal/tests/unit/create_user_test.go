package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nibroos/elearning-go/service/internal/mocks"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := &models.User{
		Username: utils.Ptr("johndoe"),
		Email:    "johndoe@example.com",
		Password: "password123",
	}

	roleIDs := []uint32{1, 2}

	hashedPassword, _ := utils.HashPassword(user.Password)
	expectedUser := &models.User{
		ID:       1,
		Username: utils.Ptr("johndoe"),
		Email:    "johndoe@example.com",
		Password: hashedPassword,
	}

	// Create a mock gorm.DB object
	mockDB, _ := gorm.Open(postgres.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	tests := []struct {
		name          string
		user          *models.User
		roleIDs       []uint32
		mockBeginTx   *gorm.DB
		mockCreateErr error
		mockAttachErr error
		mockCommitErr error
		expectedUser  *models.User
		expectedErr   error
	}{
		{
			name:          "success",
			user:          user,
			roleIDs:       roleIDs,
			mockBeginTx:   mockDB,
			mockCreateErr: nil,
			mockAttachErr: nil,
			mockCommitErr: nil,
			expectedUser:  expectedUser,
			expectedErr:   nil,
		},
		{
			name:          "create user error",
			user:          user,
			roleIDs:       roleIDs,
			mockBeginTx:   mockDB,
			mockCreateErr: errors.New("create user error"),
			mockAttachErr: nil,
			mockCommitErr: nil,
			expectedUser:  nil,
			expectedErr:   errors.New("create user error"),
		},
		{
			name:          "attach roles error",
			user:          user,
			roleIDs:       roleIDs,
			mockBeginTx:   mockDB,
			mockCreateErr: nil,
			mockAttachErr: errors.New("attach roles error"),
			mockCommitErr: nil,
			expectedUser:  nil,
			expectedErr:   errors.New("attach roles error"),
		},
		{
			name:          "commit transaction error",
			user:          user,
			roleIDs:       roleIDs,
			mockBeginTx:   mockDB,
			mockCreateErr: nil,
			mockAttachErr: nil,
			mockCommitErr: errors.New("commit transaction error"),
			expectedUser:  nil,
			expectedErr:   errors.New("commit transaction error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("BeginTransaction").Return(tt.mockBeginTx)
			mockRepo.On("CreateUser", tt.mockBeginTx, mock.AnythingOfType("*models.User")).Return(tt.mockCreateErr)
			mockRepo.On("AttachRoles", tt.mockBeginTx, mock.AnythingOfType("*models.User"), tt.roleIDs).Return(tt.mockAttachErr)
			mockRepo.On("Commit").Return(tt.mockCommitErr)
			mockRepo.On("Rollback").Return(nil)

			user, err := userService.CreateUser(ctx, tt.user, tt.roleIDs)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
			mockRepo.AssertExpectations(t)
		})
	}
}
