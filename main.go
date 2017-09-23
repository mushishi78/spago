package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v", err)
	}

	serv, err := ServerCreate(cwd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SPAGO! Listening on 127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", serv))
}
