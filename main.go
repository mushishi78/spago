package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Too many arguments provided")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v", err)
	}
	if len(os.Args) == 2 {
		cwd = filepath.Join(cwd, os.Args[1])
	}

	port := flag.Int("PORT", 8080, "the port that the dev server will listen on")
	apiPort := flag.Int("API_PORT", 3000, "the port that /api requests will be forwarded to")
	flag.Parse()

	serv, err := serverCreate(cwd, *apiPort)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listening on http://localhost:%v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), serv))
}
