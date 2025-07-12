package foodrecipe_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/klins/devpool/go-day6/wongnok/internal/foodrecipe"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		service := foodrecipe.NewService(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(service))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type ServiceCreateTestSuite struct {
	suite.Suite

	service foodrecipe.IService
	repo    *MockIRepository

	errRepositoryCreate error
}

func (suite *ServiceCreateTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.errRepositoryCreate = nil

	suite.repo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		recipe := args.Get(0).(*model.FoodRecipe)
		*recipe = model.FoodRecipe{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			CookingDurationID: 1,
			DifficultyID:      1,
		}
	}).Return(func(*model.FoodRecipe) error {
		return suite.errRepositoryCreate
	})
}

func (suite *ServiceCreateTestSuite) TestReturnRecipeCreated() {
	expectedRecipe := model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	}

	recipe, err := suite.service.Create(dto.FoodRecipeRequest{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	})
	suite.NoError(err)

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "Create", &model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	})
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRequestValidate() {
	recipe, err := suite.service.Create(dto.FoodRecipeRequest{})
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "request invalid"))

	suite.Empty(recipe)
	suite.repo.AssertNotCalled(suite.T(), "Create")
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRepositoryCreate() {
	suite.errRepositoryCreate = assert.AnError

	recipe, err := suite.service.Create(dto.FoodRecipeRequest{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	})
	suite.ErrorIs(err, assert.AnError)

	suite.Empty(recipe)
}

func TestServiceCreate(t *testing.T) {
	suite.Run(t, new(ServiceCreateTestSuite))
}
