package foodrecipe

import (
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/global"
	"github.com/klins/devpool/go-day6/wongnok/internal/helper"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	Create(request dto.FoodRecipeRequest, claims model.Claims) (model.FoodRecipe, error)
	Update(request dto.FoodRecipeRequest, id string, claims model.Claims) (model.FoodRecipe, error)
	GetByID(id string) (model.FoodRecipe, error)
	GetAll() ([]model.FoodRecipe, error)
	Get(foodRecipeQuery model.FoodRecipeQuery) (model.FoodRecipes, int64, error)
	Count() (int64, error)
	Delete(id string, claims model.Claims) error
	GetFavorites(foodRecipeQuery model.FoodRecipeQuery, claims model.Claims) (model.FoodRecipes, int64, error)
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) Create(request dto.FoodRecipeRequest, claims model.Claims) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}

	var recipe model.FoodRecipe
	recipe = recipe.FromRequest(request, claims)

	if err := service.Repository.Create(&recipe); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "create recipe")
	}

	return recipe, nil
}

func (service Service) GetByID(id string) (model.FoodRecipe, error) {
	recipe, err := service.Repository.GetByID(id)
	if err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "get recipe by ID")
	}

	// Calculate the average rating for the recipe
	recipe = helper.CalculateAverageRating(recipe)

	return recipe, nil
}

func (service Service) GetAll() ([]model.FoodRecipe, error) {
	recipes, err := service.Repository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get all recipes")
	}

	// Calculate the average rating for each recipe
	for i, recipe := range recipes {
		recipes[i] = helper.CalculateAverageRating(recipe)
	}

	return recipes, nil
}

func (service Service) Get(foodRecipeQuery model.FoodRecipeQuery) (model.FoodRecipes, int64, error) {
	total, err := service.Repository.Count()
	if err != nil {
		return nil, 0, err
	}

	results, err := service.Repository.Get(foodRecipeQuery)
	if err != nil {
		return nil, 0, err
	}

	results = helper.CalculateAverageRatings(results)

	return results, total, nil
}

func (service Service) Count() (int64, error) {
	count, err := service.Repository.Count()
	if err != nil {
		return 0, errors.Wrap(err, "count recipes")
	}
	return count, nil
}

func (service Service) CountFavorites(userID string) (int64, error) {
	count, err := service.Repository.CountFavorites(userID)
	if err != nil {
		return 0, errors.Wrap(err, "count favorite recipes")
	}
	return count, nil
}

func (service Service) Update(request dto.FoodRecipeRequest, id string, claims model.Claims) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}

	recipe, err := service.Repository.GetByID(id)
	if err != nil {
		// กรณีไม่พบ id ที่ต้องการ update
		return model.FoodRecipe{}, errors.Wrap(err, "find recipe")
	}

	if recipe.UserID != claims.ID {
		// กรณี user ที่ login ไม่ตรงกับ user ที่สร้าง recipe
		return model.FoodRecipe{}, global.ErrorForbidden
	}

	recipe = recipe.FromRequest(request, claims)

	if err := service.Repository.Update(&recipe); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "update recipe")
	}

	recipe = helper.CalculateAverageRating(recipe)

	return recipe, nil
}

func (service Service) Delete(id string, claims model.Claims) error {
	recipe, err := service.Repository.GetByID(id)
	if err != nil {
		// กรณีไม่พบ id ที่ต้องการ update
		return errors.Wrap(err, "find recipe")

	}

	if recipe.UserID != claims.ID {
		// กรณี user ที่ login ไม่ตรงกับ user ที่สร้าง recipe
		return global.ErrorForbidden
	}

	return service.Repository.Delete(id)
}

func (service Service) GetFavorites(foodRecipeQuery model.FoodRecipeQuery, claims model.Claims) (model.FoodRecipes, int64, error) {
	if claims.ID == "" {
		return nil, 0, global.ErrorForbidden
	}
	total, err := service.Repository.CountFavorites(claims.ID)
	if err != nil {
		return nil, 0, err
	}
	results, err := service.Repository.GetFavorites(foodRecipeQuery, claims.ID)
	if err != nil {
		return nil, 0, err
	}
	results = helper.CalculateAverageRatings(results)
	return results, total, nil
}
