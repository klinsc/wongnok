package user

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	GetByID(id string) (model.User, error)
	Upsert(user *model.User) error
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{DB: db}
}

func (repo Repository) GetByID(id string) (model.User, error) {
	var user model.User
	if err := repo.DB.First(&user, "id = ?", id).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (repo Repository) Upsert(user *model.User) error {
	return repo.DB.Save(user).Error
}
