// main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/klins/devpool/go-day6/wongnok/config"
	"github.com/klins/devpool/go-day6/wongnok/internal/foodrecipe"
	"github.com/klins/devpool/go-day6/wongnok/internal/rating"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}

	fmt.Println(dir)

	// Load .env from project root
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or failed to load .env")
	}

	// Load configuration
	var conf config.Config

	if err := env.Parse(&conf); err != nil {
		log.Fatal("Error when decoding configuration:", err)
	}

	// Add this line to debug the loaded configuration
	log.Printf("Attempting to connect with DSN: %s", conf.Database.URL)

	// Database connection
	db, err := gorm.Open(postgres.Open("postgres://postgres:pass2word@localhost:5432/wongnok?sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error when connect to database:", err)
	}
	// Ensure close connection when terminated
	defer func() {
		sqldb, _ := db.DB()
		sqldb.Close()
	}()

	// Handler
	foodRecipeHandler := foodrecipe.NewHandler(db)
	ratingHandler := rating.NewHandler(db)

	// Router
	router := gin.Default()

	// Register route
	group := router.Group("/api/v1")
	group.GET("/food-recipes", foodRecipeHandler.Get)
	group.GET("/food-recipes/:id", foodRecipeHandler.GetByID)
	group.POST("/food-recipes", foodRecipeHandler.Create)
	group.PUT("/food-recipes/:id", foodRecipeHandler.Update)
	group.DELETE("/food-recipes/:id", foodRecipeHandler.Delete)
	group.POST("/food-recipes/:id/ratings", ratingHandler.Create)

	if err := router.Run(); err != nil {
		log.Fatal("Server error:", err)
	}
}
