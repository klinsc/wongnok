package rating

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(rating *model.Rating) error
	GetByID(id int) (model.Ratings, error)
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
