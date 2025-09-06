package dto

type FavoriteRequest struct {
	IsFavorited *bool `json:"isFavorited" validate:"required"`
}

type FavoriteResponse struct {
	FoodRecipeID uint `json:"foodRecipeID"`
	IsFavorited  bool `json:"isFavorited"`
}

type FavoritesResponse BaseListResponse[[]FavoriteResponse]
