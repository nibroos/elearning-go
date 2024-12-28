package unit_test

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nibroos/elearning-go/service/internal/config"
	"github.com/nibroos/elearning-go/service/internal/mocks"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDeleteUser(t *testing.T) {

	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join("..", "..", "..", "..", ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Retrieve database credentials from environment variables using GetTestDatabaseURL
	configDBTest := config.GetTestDatabaseURL()

	mockDB, err := gorm.Open(postgres.Open(configDBTest), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the Gorm database: %v", err)
	}

	tests := []struct {
		name               string
		userID             uint
		mockBeginTx        *gorm.DB
		mockDeleteRolesErr error
		mockDeleteUserErr  error
		mockCommitErr      error
		expectedErr        error
	}{
		{
			name:               "success",
			userID:             1,
			mockBeginTx:        mockDB.Begin(),
			mockDeleteRolesErr: nil,
			mockDeleteUserErr:  nil,
			mockCommitErr:      nil,
			expectedErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBeginTx != nil {
				mockRepo.On("BeginTransaction").Return(tt.mockBeginTx)
				mockRepo.On("DeleteRolesByUserID", tt.mockBeginTx, tt.userID).Return(tt.mockDeleteRolesErr)
				mockRepo.On("DeleteUser", tt.mockBeginTx, tt.userID).Return(tt.mockDeleteUserErr)
				mockRepo.On("Commit", tt.mockBeginTx).Return(tt.mockCommitErr)
			} else {
				mockRepo.On("BeginTransaction").Return(nil)
			}

			err := userService.DeleteUser(ctx, tt.userID)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
