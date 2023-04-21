package main

import (
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) <= 1 {
		panic("Please run as 'go run . rest dev' or 'go run . socket dev")
	}

	var env string

	if(args[1] == "dev"){
		env = ".env.dev"
	} else if (args[1] == "prod"){
		env = ".env.prod"
	}

	if args[0] == "rest" {
		MainRestServer(env)
	} else {
		MainSocketServer(env)
	}
}