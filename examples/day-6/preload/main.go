package main

import (
	"encoding/json"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Difficulty struct {
	gorm.Model
	Name string
}

type CookingDuration struct {
	gorm.Model
	Name string
}

type FoodRecipe struct {
	gorm.Model
	Name              string
	CookingDurationID uint
	CookingDuration   CookingDuration
	DifficultyID      uint
	Difficulty        Difficulty
}

func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:pass2word@localhost:5432/wongnok?sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	var recipe FoodRecipe

	if err := db.Preload(clause.Associations).First(&recipe, 1).Error; err != nil {
		log.Fatal("failed to fetch recipe:", err)
	}

	disp, _ := json.MarshalIndent(recipe, "", "  ")
	log.Println("Fetched Recipe:", string(disp))
}
