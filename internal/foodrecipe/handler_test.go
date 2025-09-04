package foodrecipe_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/foodrecipe"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		handler := foodrecipe.NewHandler(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(handler))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type HandlerCreateTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	respServiceCreate model.FoodRecipe
	errServiceCreate  error

	// Claims
	claims model.Claims

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (suite *HandlerCreateTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerCreateTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
		Service: suite.service,
	}

	suite.claims = model.Claims{
		ID: "user-id",
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()

		// Set context
		router.Use(func(ctx *gin.Context) {
			ctx.Set("claims", suite.claims)
		})

		router.POST("/api/v1/food-recipes", suite.handler.Create)

		// Recorder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodPost,
			"/api/v1/food-recipes",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceCreate = model.FoodRecipe{
		Name: "Name",
	}
	suite.errServiceCreate = nil

	suite.service.On("Create", mock.Anything, mock.Anything).Return(func(dto.FoodRecipeRequest, model.Claims) (model.FoodRecipe, error) {
		return suite.respServiceCreate, suite.errServiceCreate
	})
}

func (suite *HandlerCreateTestSuite) TestResponseRecipeWithStatusCode201() {
	payload := strings.NewReader(`{"name":"Name"}`)

	response := suite.server(payload)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	// Expect body
	expectedBody := dto.FoodRecipeResponse{
		Name: "Name",
	}
	expectedJson, _ := json.Marshal(expectedBody)

	suite.Equal(http.StatusCreated, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
	suite.service.AssertCalled(suite.T(), "Create", dto.FoodRecipeRequest{
		Name: "Name",
	}, suite.claims)
}

func (suite *HandlerCreateTestSuite) TestErrorWhenRequestInvalid() {
	payload := strings.NewReader(``)

	response := suite.server(payload)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":"EOF"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "Create")
}

func (suite *HandlerCreateTestSuite) TestErrorWhenServiceCreateRecipe() {
	suite.errServiceCreate = assert.AnError
	payload := strings.NewReader(`{"name":"Name"}`)

	response := suite.server(payload)

	// Ensure close reader when terminated
	body := response.Result().Body
	body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func (suite *HandlerCreateTestSuite) TestValidationErrorsErrorWhenServiceCreateRecipte() {
	suite.errServiceCreate = make(validator.ValidationErrors, 0)

	payload := strings.NewReader(`{"name":"Name"}`)

	response := suite.server(payload)

	// Ensure close reader when terminated
	body := response.Result().Body
	body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":""}`, response.Body.String())
}

func TestHandlerCreate(t *testing.T) {
	suite.Run(t, new(HandlerCreateTestSuite))
}

type HandlerGetByIDTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	respRecipeInServiceGetByID model.FoodRecipe
	errServiceGetByID          error

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

func (suite *HandlerGetByIDTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerGetByIDTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
		Service: suite.service,
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/food-recipes/:id", suite.handler.GetByID)

		// Recorder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/food-recipes/1",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respRecipeInServiceGetByID = model.FoodRecipe{
		Model:       gorm.Model{ID: 1},
		Name:        "Name",
		Description: "Description",
		Ingredient:  "Ingredient",
		Instruction: "Instruction",
		CookingDuration: model.CookingDuration{
			Model: gorm.Model{ID: 1},
			Name:  "5 - 10",
		},
		Difficulty: model.Difficulty{
			Model: gorm.Model{ID: 1},
			Name:  "5 - 10",
		},
	}

	suite.errServiceGetByID = nil

	suite.service.On("GetByID", mock.AnythingOfType("string")).Return(func(id string) (model.FoodRecipe, error) {
		if id == "1" {
			return suite.respRecipeInServiceGetByID, suite.errServiceGetByID
		}
		return model.FoodRecipe{}, gorm.ErrRecordNotFound
	})
}

func (suite *HandlerGetByIDTestSuite) TestResponseRecipeWithStatusCode200() {
	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	// Expect body
	expectedBody := dto.FoodRecipeResponse{
		ID:          1,
		Name:        "Name",
		Description: "Description",
		Ingredient:  "Ingredient",
		Instruction: "Instruction",
		CookingDuration: dto.CookingDurationResponse{
			ID:   1,
			Name: "5 - 10",
		},
		Difficulty: dto.DifficultyResponse{
			ID:   1,
			Name: "5 - 10",
		},
	}
	expectedJson, _ := json.Marshal(expectedBody)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
}

func (suite *HandlerGetByIDTestSuite) TestErrorWhenRecipeNotFound() {
	suite.errServiceGetByID = gorm.ErrRecordNotFound

	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	// suite.Equal(http.StatusInternalServerError, response.Code) // passes if response.Code == 500
	suite.Equal(http.StatusNotFound, response.Code) // passes if response.Code == 404
	suite.Equal(`{"message":"Recipe not found"}`, response.Body.String())
	suite.service.AssertCalled(suite.T(), "GetByID", "1")
}

func TestHandlerGetByID(t *testing.T) {
	suite.Run(t, new(HandlerGetByIDTestSuite))
}
