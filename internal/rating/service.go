package rating

import (
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/klins/devpool/go-day6/wongnok/internal/user"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	Create(request dto.RatingRequest, recipeID int, claims model.Claims) (model.Rating, error)
	GetByID(id int) (model.Ratings, error)
	GetMyFavorites(claims model.Claims) (model.FoodRecipes, error)
	IsFavorite(recipeID int, claims model.Claims) (bool, error)
	Favorite(request dto.FavoriteRequest, recipeID int, claims model.Claims) (bool, error)
}

type Service struct {
	Repository   IRepository
	IUserService user.IService
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository:   NewRepository(db),
		IUserService: user.NewService(db),
	}
}

func (service Service) Create(request dto.RatingRequest, recipeID int, claims model.Claims) (model.Rating, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.Rating{}, errors.Wrap(err, "request invalid")
	}

	// Verify user
	user, err := service.IUserService.GetByID(claims.ID)
	if err != nil {
		return model.Rating{}, errors.Wrap(err, "get user by ID")
	}

	var rating model.Rating
	rating = rating.FromRequest(request)
	rating.FoodRecipeID = uint(recipeID)
	rating.UserID = user.ID

	if err := service.Repository.Create(&rating); err != nil {
		return model.Rating{}, errors.Wrap(err, "create recipe")
	}

	return rating, nil
}

func (service Service) GetByID(id int) (model.Ratings, error) {
	ratings, err := service.Repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (service Service) GetMyFavorites(claims model.Claims) (model.FoodRecipes, error) {
	// Verify user
	user, err := service.IUserService.GetByID(claims.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user by ID")
	}
	return service.IUserService.GetMyFavorites(user.ID)
}

func (service Service) IsFavorite(recipeID int, claims model.Claims) (bool, error) {
	// Verify user
	user, err := service.IUserService.GetByID(claims.ID)
	if err != nil {
		return false, errors.Wrap(err, "get user by ID")
	}
	return service.Repository.IsFavorite(recipeID, user.ID)
}

func (service Service) Favorite(request dto.FavoriteRequest, recipeID int, claims model.Claims) (bool, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return false, errors.Wrap(err, "request invalid")
	}
	// Verify user
	user, err := service.IUserService.GetByID(claims.ID)
	if err != nil {
		return false, errors.Wrap(err, "get user by ID")
	}

	if *request.IsFavorited {
		// Add to favorites
		return service.Repository.AddFavorite(recipeID, user.ID)
	} else {
		// Remove from favorites
		return service.Repository.RemoveFavorite(recipeID, user.ID)
	}
}
