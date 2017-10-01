package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

type server struct {
	CWD      string
	APIProxy *httputil.ReverseProxy
}

func serverCreate(cwd string, apiPort int) (*server, error) {
	u, err := url.Parse(fmt.Sprintf("http://localhost:%v", apiPort))
	if err != nil {
		return nil, fmt.Errorf("could not create api proxy url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	return &server{cwd, proxy}, nil
}

func (serv *server) listenAndServe(port int) {
	fmt.Printf("listening on http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), serv))
}

func (serv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		serv.APIProxy.ServeHTTP(w, r)
		return
	}

	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".css") ||
		strings.HasSuffix(r.URL.Path, ".js") ||
		strings.HasSuffix(r.URL.Path, ".map") ||
		strings.HasSuffix(r.URL.Path, ".png") ||
		strings.HasSuffix(r.URL.Path, ".ico") ||
		strings.HasSuffix(r.URL.Path, ".jpg") {

		http.ServeFile(w, r, filepath.Join(serv.CWD, r.URL.Path))
		return
	}

	var linkElements = make([]string, 0)
	var scriptElements = make([]string, 0)
	{
		err := unorderedWalk(serv.CWD, "", func(path string) {
			if strings.HasSuffix(path, ".css") {
				line := fmt.Sprintf("  <link href=\"/%v\" rel=\"stylesheet\" type=\"text/css\">\n  ", path)
				linkElements = append(linkElements, line)
			}

			if strings.HasSuffix(path, ".js") {
				line := fmt.Sprintf("  <script src=\"/%v\"></script>\n  ", path)
				scriptElements = append(scriptElements, line)
			}
		})
		if err != nil {
			http.Error(w, "failed scan for files", 500)
			log.Printf("failed scan for files: %v", err)
			return
		}
	}

	htmlPath := filepath.Join(serv.CWD, "index.html")

	htmlBytes, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "failed to read index.html", 500)
		log.Printf("failed to read index.html: %v", err)
		return
	}

	html := string(htmlBytes)
	cssInsertPoint := strings.Index(html, "</head>")
	jsInsertPoint := strings.Index(html, "</body>")

	if cssInsertPoint == -1 {
		http.Error(w, "index.html does not have <head> element", 500)
		return
	}
	if jsInsertPoint == -1 {
		http.Error(w, "index.html does not have <body> element", 500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, html[:cssInsertPoint])
	fmt.Fprint(w, strings.Join(linkElements, ""))
	fmt.Fprint(w, html[cssInsertPoint:jsInsertPoint])
	fmt.Fprint(w, strings.Join(scriptElements, ""))
	fmt.Fprint(w, html[jsInsertPoint:])
}