package main

import (
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) <= 0 {
		panic("Please run as 'go run . rest' or 'go run . socket")
	}

	if args[0] == "rest" {
		MainRestServer()
	} else {
		MainSocketServer()
	}
}