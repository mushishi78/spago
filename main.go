package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Too many arguments provided")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v\n", err)
	}
	if len(os.Args) == 2 {
		cwd = filepath.Join(cwd, os.Args[1])
	}

	serv, err := serverCreate(cwd, log.Printf)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(serv.listenAndServe())
}
