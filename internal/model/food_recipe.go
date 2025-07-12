package model

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"gorm.io/gorm"
)

type FoodRecipe struct {
	gorm.Model
	Name              string
	Description       string
	Ingredient        string
	Instruction       string
	ImageURL          *string
	CookingDurationID uint
	CookingDuration   CookingDuration
	DifficultyID      uint
	Difficulty        Difficulty
}

func (recipe FoodRecipe) FromRequest(request dto.FoodRecipeRequest) FoodRecipe {
	return FoodRecipe{
		Name:              request.Name,
		Description:       request.Description,
		Ingredient:        request.Ingredient,
		Instruction:       request.Instruction,
		ImageURL:          request.ImageURL,
		CookingDurationID: request.CookingDurationID,
		DifficultyID:      request.DifficultyID,
	}
}

func (recipe FoodRecipe) ToResponse() dto.FoodRecipeResponse {
	return dto.FoodRecipeResponse{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: recipe.Description,
		Ingredient:  recipe.Ingredient,
		Instruction: recipe.Instruction,
		ImageURL:    recipe.ImageURL,
		CookingDuration: dto.CookingDurationResponse{
			ID:   recipe.CookingDuration.ID,
			Name: recipe.CookingDuration.Name,
		},
		Difficulty: dto.DifficultyResponse{
			ID:   recipe.Difficulty.ID,
			Name: recipe.Difficulty.Name,
		},
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
	}
}
