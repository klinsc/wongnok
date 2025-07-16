package model

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	Score        float64
	FoodRecipeID uint
}

func (rating Rating) FromRequest(request dto.RatingRequest) Rating {
	return Rating{
		Score: request.Score,
	}
}

func (rating Rating) ToResponse() dto.RatingResponse {
	return dto.RatingResponse{
		Score:        rating.Score,
		FoodRecipeID: rating.FoodRecipeID,
	}
}
