package unit_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nibroos/elearning-go/service/internal/mocks"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Failure in HashPassword
// Failure in Transaction Begin
// Failure in Create User
// Failure in Attach Roles
// Failure in Commit

func TestCreateUser(t *testing.T) {

	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join("..", "..", "..", "..", ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := models.User{
		Name:     "TName",
		Username: utils.Ptr("johndoe"),
		Email:    "johndoe@example.com",
		Password: "password123",
		Address:  utils.Ptr("TAddress"),
	}

	originPassword := user.Password
	hashedPassword, _ := utils.HashPassword(user.Password)

	roleIDs := []uint32{1, 2}

	expectedUser := models.User{
		ID:       0,
		Name:     "TName",
		Username: utils.Ptr("johndoe"),
		Email:    "johndoe@example.com",
		Password: hashedPassword,
		Address:  utils.Ptr("TAddress"),
	}

	// Retrieve database credentials from environment variables
	configDBTest := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB_TEST"),
	)

	// mockDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	mockDB, err := gorm.Open(postgres.Open(configDBTest), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the Gorm database: %v", err)
	}

	tests := []struct {
		name          string
		user          models.User
		roleIDs       []uint32
		mockBeginTx   *gorm.DB
		mockCreateErr error
		mockAttachErr error
		mockCommitErr error
		expectedUser  *models.User
		expectedErr   error
	}{
		{
			name:          "success create user",
			user:          user,
			roleIDs:       roleIDs,
			mockBeginTx:   mockDB.Begin(),
			mockCreateErr: nil,
			mockAttachErr: nil,
			mockCommitErr: nil,
			expectedUser:  &expectedUser,
			expectedErr:   nil,
		},
		{
			name:          "error hash password",
			user:          models.User{Password: ""},
			roleIDs:       roleIDs,
			mockBeginTx:   nil,
			mockCreateErr: nil,
			mockAttachErr: nil,
			mockCommitErr: nil,
			expectedUser:  nil,
			expectedErr:   bcrypt.ErrHashTooShort,
		},
		{
			name:          "error roleIDs empty error",
			user:          user,
			roleIDs:       []uint32{},
			mockBeginTx:   nil,
			mockCreateErr: nil,
			mockAttachErr: nil,
			mockCommitErr: nil,
			expectedUser:  nil,
			expectedErr:   fmt.Errorf("roleIDs cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBeginTx != nil {
				mockRepo.On("BeginTransaction").Return(tt.mockBeginTx)
				mockRepo.On("CreateUser", tt.mockBeginTx, mock.AnythingOfType("*models.User")).Return(tt.mockCreateErr)
				mockRepo.On("AttachRoles", tt.mockBeginTx, mock.AnythingOfType("*models.User"), tt.roleIDs).Return(tt.mockAttachErr)
				mockRepo.On("Commit", tt.mockBeginTx).Return(tt.mockCommitErr)
			}

			user, err := userService.CreateUser(ctx, &tt.user, tt.roleIDs)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedUser != nil {
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Address, user.Address)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(tt.user.Password), []byte(originPassword)))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
