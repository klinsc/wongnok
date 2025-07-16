package rating

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(rating *model.Rating) error
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
