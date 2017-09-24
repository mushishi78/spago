package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
		strings.HasSuffix(r.URL.Path, ".png") ||
		strings.HasSuffix(r.URL.Path, ".ico") ||
		strings.HasSuffix(r.URL.Path, ".jpg") {

		http.ServeFile(w, r, filepath.Join(serv.CWD, r.URL.Path))
		return
	}

	var linkElements = make([]string, 0)
	var scriptElements = make([]string, 0)
	{
		err := filepath.Walk(serv.CWD, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("file walk error: %v", err)
			}

			relPath := path[len(serv.CWD):]
			relPath = strings.Replace(relPath, "\\", "/", -1)

			if strings.HasSuffix(path, ".css") {
				line := fmt.Sprintf("  <link href=\"%v\" rel=\"stylesheet\" type=\"text/css\">\n  ", relPath)
				linkElements = append(linkElements, line)
				return nil
			}

			if strings.HasSuffix(path, ".js") {
				line := fmt.Sprintf("  <script src=\"%v\"></script>\n  ", relPath)
				scriptElements = append(scriptElements, line)
				return nil
			}

			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// Read index file
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
