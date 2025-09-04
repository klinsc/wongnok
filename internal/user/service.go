package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/helper"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	UpsertWithClaims(claims model.Claims) (model.User, error)
	GetByID(id string) (model.User, error)
	GetRecipes(userID string, claims model.Claims) (model.FoodRecipes, error)
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) UpsertWithClaims(claims model.Claims) (model.User, error) {
	validate := validator.New()
	if err := validate.Struct(claims); err != nil {
		return model.User{}, err
	}

	// Get user
	user, err := service.Repository.GetByID(claims.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, errors.Wrap(err, "get user by ID")
	}

	// Set claimed information to user
	user = user.FromClaims(claims)

	// Upsert user
	if err := service.Repository.Upsert(&user); err != nil {
		return model.User{}, errors.Wrap(err, "upsert user")
	}

	return user, nil
}

func (service Service) GetByID(id string) (model.User, error) {
	user, err := service.Repository.GetByID(id)
	if err != nil {
		return model.User{}, errors.Wrap(err, "get user by ID")
	}

	return user, nil
}

func (service Service) GetRecipes(userID string, claims model.Claims) (model.FoodRecipes, error) {
	if _, err := service.Repository.GetByID(claims.ID); err != nil {
		return model.FoodRecipes{}, errors.Wrap(err, "find user")
	}

	foodRecipes, err := service.Repository.GetRecipes(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.FoodRecipes{}, errors.Wrap(err, "get recipes")
	}

	foodRecipes = helper.CalculateAverageRatings(foodRecipes)

	return foodRecipes, nil
}
