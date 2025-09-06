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
	Get(query model.FoodRecipeQuery) (model.FoodRecipes, error)
	GetFavorites(query model.FoodRecipeQuery, userID string) (model.FoodRecipes, error)
	Count() (int64, error)
	CountFavorites(userID string) (int64, error)
	Update(recipe *model.FoodRecipe) error
	Delete(id string) error
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

func (repo Repository) Get(query model.FoodRecipeQuery) (model.FoodRecipes, error) {
	var recipes = make(model.FoodRecipes, 0)

	offset := (query.Page - 1) * query.Limit
	db := repo.DB.Preload(clause.Associations)

	if query.Search != "" {
		db = db.Where("name LIKE ?", "%"+query.Search+"%").Or("description LIKE ?", "%"+query.Search+"%")
	}

	if err := db.Order("name asc").Limit(query.Limit).Offset(offset).Find(&recipes).Error; err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo Repository) GetFavorites(query model.FoodRecipeQuery, userID string) (model.FoodRecipes, error) {
	var recipes = make(model.FoodRecipes, 0)

	offset := (query.Page - 1) * query.Limit
	db := repo.DB.Preload(clause.Associations).
		Joins("JOIN favorites ON favorites.food_recipe_id = food_recipes.id").
		Where("favorites.user_id = ? AND favorites.deleted_at IS NULL", userID)
	if query.Search != "" {
		db = db.Where("food_recipes.name LIKE ?", "%"+query.Search+"%").Or("food_recipes.description LIKE ?", "%"+query.Search+"%")
	}
	if err := db.Order("food_recipes.name asc").Limit(query.Limit).Offset(offset).Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (repo Repository) Count() (int64, error) {
	var count int64
	err := repo.DB.Model(&model.FoodRecipe{}).Count(&count).Error
	return count, err
}

func (repo Repository) CountFavorites(userID string) (int64, error) {
	var count int64

	// Get count of favorite recipes for the user with null deleted_at
	err := repo.DB.Model(&model.FoodRecipe{}).
		Joins("JOIN favorites ON favorites.food_recipe_id = food_recipes.id").
		Where("favorites.user_id = ? AND favorites.deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (repo Repository) Update(recipe *model.FoodRecipe) error {
	// update
	if err := repo.DB.Model(&recipe).Updates(recipe).Error; err != nil {
		return err
	}

	return repo.DB.Preload(clause.Associations).First(&recipe, recipe.ID).Error
}

func (repo Repository) Delete(id string) error {
	return repo.DB.Delete(&model.FoodRecipes{}, id).Error
}
