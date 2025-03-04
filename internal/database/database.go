package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"blockstracker_backend/config"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		dbUser, dbPassword, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), config.GormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB = db
	fmt.Println("Connected to the database!")
}

func DBProvider() *gorm.DB {
	return DB
}
