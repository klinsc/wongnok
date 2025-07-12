package user_test

import (
	"github.com/klins/devpool/go-day6/wongnok/examples/day-6/user"
	"github.com/stretchr/testify/mock"
)

type MockIRepository struct {
	mock.Mock
}

func (m *MockIRepository) Get(id string) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockIRepository) Create(user user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockIRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
