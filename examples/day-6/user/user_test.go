package user_test

import (
	"testing"

	"github.com/klins/devpool/go-day6/wongnok/examples/day-6/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGet(t *testing.T) {
	t.Run("Should return user when found", func(t *testing.T) {
		mockRepo := new(MockIRepository)
		service := user.NewService(mockRepo)

		mockRepo.On("Get", mock.Anything).Return(user.User{Name: "Peter"}, nil)

		result, err := service.Get("id")
		assert.NoError(t, err, "Get should not return an error")

		assert.Equal(t, user.User{Name: "Peter"}, result, "Expected user should be Peter")
	})

	t.Run("Should return error when user not found", func(t *testing.T) {
		mockRepo := new(MockIRepository)
		service := user.NewService(mockRepo)

		mockRepo.On("Get", mock.Anything).Return(user.User{}, assert.AnError)

		result, err := service.Get("id")
		assert.Error(t, err, "Get should return an error when user not found")

		assert.Equal(t, user.User{}, result, "Expected result should be an empty user")
	})
}
