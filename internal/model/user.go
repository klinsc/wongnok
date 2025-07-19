package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string
	FirstName string
	LastName  string
}

func (user User) FromClaims(claims Claims) User {
	return User{
		Model:     user.Model,
		ID:        claims.ID,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
	}
}
