package utils

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	dbHost := os.Getenv("dbHost")
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	dbName := os.Getenv("dbName")
	dbPort := os.Getenv("dbPort")
	dbSSLMode := os.Getenv("dbSSLMode")

	dsn := fmt.Sprintf("host= %s user= %s dbname= %s port= %s sslmode= %s", dbHost, dbUser, dbName, dbPort, dbSSLMode)
	if(len([]byte(dbPass)) > 0){
		dsn = dsn + fmt.Sprintf(" password= %s", dbPass)
	}

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
