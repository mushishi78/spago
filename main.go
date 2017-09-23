package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type server struct {
	CWD string
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v", err)
	}

	port := flag.Int("PORT", 8080, "the port that the dev server will listen on")
	flag.Parse()

	serv := &server{cwd}
	fmt.Printf("listening on http://localhost:%v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), serv))
}
