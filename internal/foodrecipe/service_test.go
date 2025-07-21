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
	claims := model.Claims{
		ID: "user-id",
	}

	expectedRecipe := model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "user-id", // new, set the user ID from claims
	}

	recipe, err := suite.service.Create(dto.FoodRecipeRequest{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	}, claims)
	suite.NoError(err)

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "Create", &model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "user-id", // new, set the user ID from claims
	})
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRequestValidate() {
	claims := model.Claims{
		ID: "user-id",
	}

	recipe, err := suite.service.Create(dto.FoodRecipeRequest{}, claims)
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "request invalid"))

	suite.Empty(recipe)
	suite.repo.AssertNotCalled(suite.T(), "Create")
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRepositoryCreate() {
	claims := model.Claims{
		ID: "user-id",
	}

	suite.errRepositoryCreate = assert.AnError

	request := dto.FoodRecipeRequest{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	}
	recipe, err := suite.service.Create(request, claims)
	suite.Error(err)
	suite.EqualError(err, "create recipe: "+assert.AnError.Error())

	suite.Empty(recipe)
	suite.repo.AssertCalled(suite.T(), "Create", &model.FoodRecipe{
		Name:              request.Name,
		Description:       request.Description,
		Ingredient:        request.Ingredient,
		Instruction:       request.Instruction,
		CookingDurationID: request.CookingDurationID,
		DifficultyID:      request.DifficultyID,
	})
	suite.repo.AssertExpectations(suite.T())

}

func TestServiceCreate(t *testing.T) {
	suite.Run(t, new(ServiceCreateTestSuite))
}

type ServiceGetByIDTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	errRepositoryGetByID      error
	responseRepositoryGetByID model.FoodRecipe
}

func (suite *ServiceGetByIDTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.errRepositoryGetByID = nil
	suite.responseRepositoryGetByID = model.FoodRecipe{
		Name: "Name",
	}

	suite.repo.On("GetByID", mock.AnythingOfType("string")).Return(func(id string) (model.FoodRecipe, error) {
		if id == "1" {
			return suite.responseRepositoryGetByID, suite.errRepositoryGetByID
		}
		return model.FoodRecipe{}, gorm.ErrRecordNotFound
	})
}

func (suite *ServiceGetByIDTestSuite) TestReturnRecipeWhenFound() {
	recipe, err := suite.service.GetByID("1")
	suite.NoError(err)

	expectedRecipe := model.FoodRecipe{
		Name: "Name",
	}

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "GetByID", "1")
}

func (suite *ServiceGetByIDTestSuite) TestErrorWhenRecipeNotFound() {
	recipe, err := suite.service.GetByID("2")
	suite.ErrorIs(err, gorm.ErrRecordNotFound)
	suite.Empty(recipe)
	suite.repo.AssertCalled(suite.T(), "GetByID", "2")
}

func TestServiceGetByID(t *testing.T) {
	suite.Run(t, new(ServiceGetByIDTestSuite))
}
