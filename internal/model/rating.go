package model

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	Score        float64
	FoodRecipeID uint
	UserID       string
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

type Ratings []Rating

func (ratings Ratings) ToResponse() dto.RatingsResponse {
	var results = make([]dto.RatingResponse, 0)

	for _, rating := range ratings {
		results = append(results, rating.ToResponse())
	}

	return dto.RatingsResponse{
		Results: results,
	}
}
