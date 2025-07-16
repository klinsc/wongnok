package foodrecipe_test

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/klins/devpool/go-day6/wongnok/internal/foodrecipe"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewRepository(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		repo := foodrecipe.NewRepository(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(repo))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type RepositoryTestSuite struct {
	suite.Suite
	ctx       context.Context
	container *postgres.PostgresContainer
	db        *gorm.DB
	repo      foodrecipe.IRepository
}

func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	path := filepath.Join("..", "..", "tests", "init-db.sql")
	fmt.Println("Using init script at:", path)

	container, err := postgres.Run(
		suite.ctx,
		"postgres:17-alpine",
		postgres.WithInitScripts(path),
		postgres.WithDatabase("wongnok-test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(
				(5 * time.Second),
			),
		),
	)
	suite.NoError(err)
	suite.container = container
}

// This will run once after all tests in the suite
func (suite *RepositoryTestSuite) TearDownSuite() {
	err := suite.container.Terminate(suite.ctx)
	suite.NoError(err)
}

// This will run before each test
func (suite *RepositoryTestSuite) SetupTest() {
	conn, err := suite.container.ConnectionString(suite.ctx, "sslmode=disable")
	suite.NoError(err)

	db, err := gorm.Open(driver.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	suite.NoError(err)

	suite.repo = &foodrecipe.Repository{
		DB: db,
	}

	suite.db = db
}

// This will run after each test
func (suite *RepositoryTestSuite) TearDownTest() {
	sqldb, _ := suite.db.DB()
	sqldb.Close()
}

// Extend
type RepositoryCreateTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryCreateTestSuite) TestReturnRecipePointerWhenCreated() {
	recipe := model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	}

	err := suite.repo.Create(&recipe)
	suite.NoError(err)

	expectedRecipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 2},
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		CookingDuration: model.CookingDuration{
			Model: gorm.Model{ID: 1},
			Name:  "5 - 10",
		},
		DifficultyID: 1,
		Difficulty: model.Difficulty{
			Model: gorm.Model{ID: 1},
			Name:  "Easy",
		},
	}

	recipe.CookingDuration.CreatedAt, recipe.CookingDuration.UpdatedAt = time.Time{}, time.Time{}
	recipe.Difficulty.CreatedAt, recipe.Difficulty.UpdatedAt = time.Time{}, time.Time{}
	recipe.CreatedAt, recipe.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedRecipe, recipe)
}

func (suite *RepositoryCreateTestSuite) TestErrorWhenDuplicateID() {
	recipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
	}

	err := suite.repo.Create(&recipe)
	suite.Error(err)
	suite.IsType(err, &pgconn.PgError{})

	suite.Equal("23505", err.(*pgconn.PgError).SQLState())
}

func TestRepositoryCreate(t *testing.T) {
	suite.Run(t, new(RepositoryCreateTestSuite))
}

type RepositoryGetByIDTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryGetByIDTestSuite) TestReturnRecipeByID() {
	recipeID := "1"

	recipe, err := suite.repo.GetByID(recipeID)
	suite.NoError(err)

	expectedRecipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "Omlet",
		Description:       "Eggs fried?",
		Ingredient:        "Eggs",
		Instruction:       "Cooking",
		CookingDurationID: 1,
		CookingDuration: model.CookingDuration{
			Model: gorm.Model{ID: 1},
			Name:  "5 - 10",
		},
		DifficultyID: 1,
		Difficulty: model.Difficulty{
			Model: gorm.Model{ID: 1},
			Name:  "Easy",
		},
	}

	// Ignore CreatedAt and UpdatedAt fields
	recipe.CookingDuration.CreatedAt, recipe.CookingDuration.UpdatedAt = time.Time{}, time.Time{}
	recipe.Difficulty.CreatedAt, recipe.Difficulty.UpdatedAt = time.Time{}, time.Time{}
	recipe.CreatedAt, recipe.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedRecipe, recipe)
}

func (suite *RepositoryGetByIDTestSuite) TestErrorWhenRecipeNotFound() {
	recipeID := "999" // Assuming this ID does not exist

	recipe, err := suite.repo.GetByID(recipeID)
	suite.Error(err)
	suite.EqualError(err, "record not found")

	suite.Empty(recipe)
}

func TestRepositoryGetByID(t *testing.T) {
	suite.Run(t, new(RepositoryGetByIDTestSuite))
}
