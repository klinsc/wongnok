package foodrecipe

import (
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	Create(request dto.FoodRecipeRequest) (model.FoodRecipe, error)
	GetByID(id string) (model.FoodRecipe, error)
	GetAll() ([]model.FoodRecipe, error)
	Get() (model.FoodRecipes, int64, error)
	Count() (int64, error)
	Update(id string, request dto.FoodRecipeRequest) (model.FoodRecipe, error)
	Delete(id string) error
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) Create(request dto.FoodRecipeRequest) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}

	var recipe model.FoodRecipe
	recipe = recipe.FromRequest(request)

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

	return recipe, nil
}

func (service Service) GetAll() ([]model.FoodRecipe, error) {
	recipes, err := service.Repository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get all recipes")
	}

	return recipes, nil
}

func (service Service) Get() (model.FoodRecipes, int64, error) {
	// Call the repository method to get recipes and total count
	recipes, err := service.Repository.Get()
	if err != nil {
		return nil, 0, errors.Wrap(err, "get recipes")
	}

	// Return the recipes and get the total count
	total, err := service.Repository.Count()
	if err != nil {
		return nil, 0, errors.Wrap(err, "count recipes")
	}

	return recipes, total, nil
}

func (service Service) Count() (int64, error) {
	count, err := service.Repository.Count()
	if err != nil {
		return 0, errors.Wrap(err, "count recipes")
	}
	return count, nil
}

func (service Service) Update(id string, request dto.FoodRecipeRequest) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}
	var recipe model.FoodRecipe

	recipe = recipe.FromRequest(request)
	if err := service.Repository.Update(id, &recipe); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "update recipe")
	}
	updatedRecipe, err := service.Repository.GetByID(id)
	if err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "get updated recipe by ID")
	}
	return updatedRecipe, nil
}

func (service Service) Delete(id string) error {
	if err := service.Repository.Delete(id); err != nil {
		return errors.Wrap(err, "delete recipe")
	}
	return nil
}
