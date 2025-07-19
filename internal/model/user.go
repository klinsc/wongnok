package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string
	FirstName string
	LastName  string
}
