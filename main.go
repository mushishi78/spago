package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type server struct {
	CWD      string
	APIProxy *httputil.ReverseProxy
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v", err)
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

func serverCreate(cwd string, apiPort int) (*server, error) {
	u, err := url.Parse(fmt.Sprintf("http://localhost:%v", apiPort))
	if err != nil {
		return nil, fmt.Errorf("could not create api proxy url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	return &server{cwd, proxy}, nil
}
