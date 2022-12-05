package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	dsn := "host=localhost user=postgres password=toor dbname=rt port=5432 sslmode=disable"

	var dbErr error

	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// DB Connection
	if dbErr != nil {
		fmt.Println("Could not connect to DB:", dbErr)
		return dbErr
	} else {
		fmt.Println("Connected to DB.")
		return nil
	}
}
