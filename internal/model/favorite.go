package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	FoodRecipeID uint
	UserID       string
}

type Favorites []Favorite

func (favorites Favorites) ToResponse() []uint {
	var results []uint
	for _, favorite := range favorites {
		results = append(results, favorite.FoodRecipeID)
	}
	return results
}
