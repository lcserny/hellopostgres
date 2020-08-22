package main

import (
	"hellopostgres"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("Please provide path to configs folder as first argument")
	}

	log.Fatal(hellopostgres.Run(args[0]))
}
