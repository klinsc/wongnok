package assertion

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestBasicAssertion(t *testing.T) {
	t.Log("Equal and NotEqual assertions")
	assert.Equal(t, 1, 1, "they should be equal")
	assert.NotEqual(t, 1, 2, "they should not be equal")

	t.Log("True and False assertions")
	assert.True(t, true, "this should be true")
	assert.False(t, false, "this should be false")

	t.Log("Nil and NotNil assertions")
	assert.Nil(t, nil, "this should be nil")
	assert.NotNil(t, "not nil", "this should not be nil")
}

func TestWithError(t *testing.T) {
	t.Log("Error assertions")
	err := func() error { return nil }()
	assert.NoError(t, err, "this should not return an error")

	err = func() error { return assert.AnError }() // using a predefined error from testify
	assert.Error(t, err, "this should return an error")

	t.Log("Empty and NotEmpty assertions")
	assert.Empty(t, "", "this should be empty")
	assert.NotEmpty(t, "not empty", "this should not be empty")
}

// Suit //
type UserTestSuite struct {
	suite.Suite // Embedding suite.Suite to create a test suite
	name        string
}

func (s *UserTestSuite) SetupSuite() {
	fmt.Println("Before all")
}

func (s *UserTestSuite) TearDownSuite() {
	fmt.Println("After all")
}

func (s *UserTestSuite) SetupTest() {
	fmt.Println("Before each test")
	s.name = "Peter"
}

func (s *UserTestSuite) TearDownTest() {
	fmt.Println("After each test")
	s.name = ""
}

func (s *UserTestSuite) TestGetUserName() {
	fmt.Println("Testing GetUserName")
	assert.Equal(s.T(), "Peter", s.name, "User name should be Peter")
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite)) // Running the test suite
}
