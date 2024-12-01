package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/mocks"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := &dtos.GetUserByIDParams{ID: 1}
	expectedUser := &dtos.UserDetailDTO{ID: 1, Name: "John Doe"}

	tests := []struct {
		name         string
		params       *dtos.GetUserByIDParams
		mockResp     *dtos.UserDetailDTO
		mockErr      error
		expectedUser *dtos.UserDetailDTO
		expectedErr  error
	}{
		{
			name:         "success",
			params:       params,
			mockResp:     expectedUser,
			mockErr:      nil,
			expectedUser: expectedUser,
			expectedErr:  nil,
		},
		{
			name:         "repository error",
			params:       params,
			mockResp:     nil,
			mockErr:      errors.New("repository error"),
			expectedUser: nil,
			expectedErr:  errors.New("repository error"),
		},
		{
			name:         "context timeout",
			params:       params,
			mockResp:     nil,
			mockErr:      context.DeadlineExceeded,
			expectedUser: nil,
			expectedErr:  context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUserByID", ctx, tt.params).Return(tt.mockResp, tt.mockErr).Once()

			user, err := userService.GetUserByID(ctx, tt.params)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
			mockRepo.AssertExpectations(t)
		})
	}
}
