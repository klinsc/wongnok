package dto

import (
	"time"
	// "github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
)

type FoodRecipeRequest struct {
	Name              string
	Description       string
	Ingredient        string
	Instruction       string
	ImageURL          *string
	CookingDurationID uint
	DifficultyID      uint
}

type FoodRecipeResponse struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Ingredient      string                  `json:"ingredient"`
	Instruction     string                  `json:"instruction"`
	ImageURL        *string                 `json:"imageUrl,omitempty"`
	CookingDuration CookingDurationResponse `json:"cookingDuration"`
	Difficulty      DifficultyResponse      `json:"difficulty"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

// type FoodRecipesResponse dto.BaseListResponse[[]FoodRecipeResponse]
type FoodRecipesResponse BaseListResponse[[]FoodRecipeResponse]
