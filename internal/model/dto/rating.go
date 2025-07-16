package dto

type RatingRequest struct {
	Score float64 `validate:"required"`
}

type RatingResponse struct {
	Score        float64 `json:"score"`
	FoodRecipeID uint    `json:"foodRecipeID"`
}

type RatingsResponse BaseListResponse[[]RatingResponse]
