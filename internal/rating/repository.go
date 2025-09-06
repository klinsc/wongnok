package rating

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository interface {
	Create(rating *model.Rating) error
	GetByID(id int) (model.Ratings, error)
	IsFavorite(recipeID int, userID string) (bool, error)
	AddFavorite(recipeID int, userID string) (bool, error)
	RemoveFavorite(recipeID int, userID string) (bool, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		DB: db,
	}
}

func (repo Repository) Create(rating *model.Rating) error {
	if err := repo.DB.Create(rating).First(&rating).Error; err != nil {
		return err
	}

	return nil
}

func (repo Repository) GetByID(id int) (model.Ratings, error) {
	var ratings model.Ratings

	if err := repo.DB.Where("food_recipe_id = ?", id).Find(&ratings).Error; err != nil {
		return nil, err
	}

	return ratings, nil
}

func (repo Repository) IsFavorite(recipeID int, userID string) (bool, error) {
	var favorite model.Favorite
	err := repo.DB.Where("food_recipe_id = ? AND user_id = ?", recipeID, userID).First(&favorite).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.Wrap(err, "query favorite")
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

func (repo Repository) AddFavorite(recipeID int, userID string) (bool, error) {
	favorite := model.Favorite{
		FoodRecipeID: uint(recipeID),
		UserID:       userID,
	}
	// Upsert favorite
	if err := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "food_recipe_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "deleted_at"}),
	}).Create(&favorite).Error; err != nil {
		return false, errors.Wrap(err, "upsert favorite")
	}

	return true, nil
}

func (repo Repository) RemoveFavorite(recipeID int, userID string) (bool, error) {
	if err := repo.DB.Where("food_recipe_id = ? AND user_id = ?", recipeID, userID).Delete(&model.Favorite{}).Error; err != nil {
		return false, errors.Wrap(err, "delete favorite")
	}
	return true, nil
}
