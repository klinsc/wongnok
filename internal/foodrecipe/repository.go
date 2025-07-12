package foodrecipe

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository interface {
	Create(recipe *model.FoodRecipe) error
	GetByID(id string) (model.FoodRecipe, error)
	GetAll() ([]model.FoodRecipe, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		DB: db,
	}
}

func (repo Repository) Create(recipe *model.FoodRecipe) error {
	return repo.DB.Preload(clause.Associations).Create(recipe).First(&recipe).Error
}

func (repo Repository) GetByID(id string) (model.FoodRecipe, error) {
	var recipe model.FoodRecipe
	err := repo.DB.Preload(clause.Associations).First(&recipe, "id = ?", id).Error
	return recipe, err
}

func (repo Repository) GetAll() ([]model.FoodRecipe, error) {
	var recipes []model.FoodRecipe
	err := repo.DB.Preload(clause.Associations).Find(&recipes).Error
	return recipes, err
}
