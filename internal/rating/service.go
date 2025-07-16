package rating

import (
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	Create(request dto.RatingRequest, recipeID int) (model.Rating, error)
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) Create(request dto.RatingRequest, recipeID int) (model.Rating, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.Rating{}, errors.Wrap(err, "request invalid")
	}

	var rating model.Rating
	rating = rating.FromRequest(request)
	rating.FoodRecipeID = uint(recipeID)

	if err := service.Repository.Create(&rating); err != nil {
		return model.Rating{}, errors.Wrap(err, "create recipe")
	}

	return rating, nil
}
