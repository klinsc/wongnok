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
	Get() (model.FoodRecipes, error)
	Count() (int64, error)
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

func (repo Repository) Get() (model.FoodRecipes, error) {
	var recipes = make(model.FoodRecipes, 0)

	err := repo.DB.Preload(clause.Associations).Find(&recipes).Error
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo Repository) Count() (int64, error) {
	var count int64
	err := repo.DB.Model(&model.FoodRecipe{}).Count(&count).Error
	return count, err
}
