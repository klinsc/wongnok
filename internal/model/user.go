package model

import (
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string
	FirstName string
	LastName  string
	ImageURL  string
}

func (user User) FromClaims(claims Claims) User {
	return User{
		Model:     user.Model,
		ID:        claims.ID,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
	}
}

func (user User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ImageURL:  user.ImageURL,
	}
}

func (user User) FromRequest(request dto.UserRequest) User {
	return User{
		Model:     user.Model,
		ID:        user.ID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		ImageURL:  request.ImageURL,
	}
}
