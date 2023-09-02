package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	// Define the PostgreSQL connection string.
	dsn := "host=localhost user=gram password=c0ld dbname=globolist port=5432 sslmode=disable TimeZone=Europe/London"

	// Open a connection to the PostgreSQL database.
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the database to create tables based on your GORM models.
	if err := DB.AutoMigrate(&Company{}, &Address{}, &PreviousName{}, &SICCode{}); err != nil {
		log.Fatalf("Failed to Automigrate: %d", err)
	}

	fmt.Println("Database migration successful")

	// Now you can use the 'db' connection to perform database operations with your GORM models.
	// For example:
	// - Creating new records
	// - Querying records
	// - Updating records
	// - Deleting records
}
