package helper

import "github.com/klins/devpool/go-day6/wongnok/internal/model"

func CalculateAverageRating(recipe model.FoodRecipe) model.FoodRecipe {
	if len(recipe.Ratings) > 0 {
		var totalRating float64
		for _, rating := range recipe.Ratings {
			totalRating += rating.Score
		}
		recipe.AverageRating = totalRating / float64(len(recipe.Ratings))
	} else {
		recipe.AverageRating = 0
	}
	return recipe
}

func CalculateAverageRatings(recipes model.FoodRecipes) model.FoodRecipes {
	for i, recipe := range recipes {
		if len(recipe.Ratings) > 0 {
			var totalRating float64
			for _, rating := range recipe.Ratings {
				totalRating += rating.Score
			}
			recipes[i].AverageRating = totalRating / float64(len(recipe.Ratings))
		} else {
			recipes[i].AverageRating = 0
		}
	}
	return recipes
}
