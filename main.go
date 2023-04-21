package main

import (
	"io/ioutil"
	"kozo/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	args := os.Args[1:]

	if len(args) <= 1 {
		panic("Please run as 'go run . rest dev' or 'go run . socket dev")
	}

	dir, _ := ioutil.ReadFile("dir.txt")
	err := godotenv.Load(".env.path")

	var env string = string(dir)

	if(args[1] == "dev"){
		env = env + ".env.dev"
	} else if (args[1] == "prod"){
		env = env + ".env.prod"
	}

	// Load Env
	err = godotenv.Load(env)

	if err != nil {
		panic(err)
	}

	// DB Connection
	dbError := utils.DBConnect()

	if dbError != nil {
		panic(dbError)
	}

	// Run migrations
	utils.DBMigrate(false)

	if args[0] == "rest" {
		MainRestServer()
	} else {
		MainSocketServer()
	}
}