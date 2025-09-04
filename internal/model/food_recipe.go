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
	Ratings           Ratings // new
	AverageRating     float64 `gorm:"-"` // new
	UserID            string  // new, user who created the recipe
	User              User    // new, relationship to User
}

type FoodRecipeQuery struct {
	Search string `form:"search"`
	Page   int    `form:"page" binding:"required,min=1"`  // page number for pagination
	Limit  int    `form:"limit" binding:"required,min=1"` // number of items per page
}

func (recipe FoodRecipe) FromRequest(request dto.FoodRecipeRequest, claims Claims) FoodRecipe {
	return FoodRecipe{
		Name:              request.Name,
		Description:       request.Description,
		Ingredient:        request.Ingredient,
		Instruction:       request.Instruction,
		ImageURL:          request.ImageURL,
		CookingDurationID: request.CookingDurationID,
		DifficultyID:      request.DifficultyID,
		UserID:            claims.ID, // new, set the user ID from claims
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
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
		AverageRating: recipe.AverageRating, // new

		User: recipe.User.ToResponse(), // new, user who created the recipe
	}
}

type FoodRecipes []FoodRecipe

func (recipes FoodRecipes) ToResponse(
	total int64,
) dto.FoodRecipesResponse {
	var result = make([]dto.FoodRecipeResponse, 0)

	for _, recipe := range recipes {
		result = append(result, recipe.ToResponse())
	}

	return dto.FoodRecipesResponse{
		// Total:   int64(len(recipes)),
		Total:   total,
		Results: result,
	}
}
