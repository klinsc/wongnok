package model_test

import (
	"testing"
	"time"

	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFoodRecipeFromRequest(t *testing.T) {
	claims := model.Claims{
		ID: "user-id",
	}

	t.Run("ShouldSetFoodRecipeModelFromRequest", func(t *testing.T) {
		imageURL := "http://example.com/image.jpg"

		request := dto.FoodRecipeRequest{
			Name:              "Test Recipe",
			Description:       "Test Description",
			Ingredient:        "Test Ingredient",
			Instruction:       "Test Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      2,
		}

		recipe := model.FoodRecipe{}.FromRequest(request, claims)

		expected := model.FoodRecipe{
			Name:              "Test Recipe",
			Description:       "Test Description",
			Ingredient:        "Test Ingredient",
			Instruction:       "Test Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      2,
		}

		assert.Equal(t, expected, recipe, "FoodRecipe should match expected values")
	})
}

func TestFoodRecipeToResponse(t *testing.T) {
	t.Run("ShouldReturnFoodRecipeResponse", func(t *testing.T) {
		imageURL := "http://example.com/image.jpg"
		mockTime := time.Now()

		recipe := model.FoodRecipe{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
			Name:        "Test Recipe",
			Description: "Test Description",
			Ingredient:  "Test Ingredient",
			Instruction: "Test Instruction",
			ImageURL:    &imageURL,
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{
					ID: 1,
				},
				Name: "Test Duration",
			},
			Difficulty: model.Difficulty{
				Model: gorm.Model{
					ID: 2,
				},
				Name: "Test Difficulty",
			},
		}
		response := recipe.ToResponse()
		expected := dto.FoodRecipeResponse{
			ID:          1,
			Name:        "Test Recipe",
			Description: "Test Description",
			Ingredient:  "Test Ingredient",

			Instruction: "Test Instruction",
			ImageURL:    &imageURL,
			CookingDuration: dto.CookingDurationResponse{
				ID:   1,
				Name: "Test Duration",
			},
			Difficulty: dto.DifficultyResponse{
				ID:   2,
				Name: "Test Difficulty",
			},
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		}
		assert.Equal(t, expected, response, "FoodRecipeResponse should match expected values")
	})
}
