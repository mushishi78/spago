package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type server struct {
	RootDir              string
	Port                 int
	ExcludedPaths        map[string]bool
	StaticFileExtensions []string
	ReverseProxyRoute    string
	ReverseProxy         *httputil.ReverseProxy
}

func serverCreate(rootDir string) (*server, error) {
	// Create server with defaults
	serv := &server{
		RootDir:              rootDir,
		Port:                 8080,
		ExcludedPaths:        map[string]bool{"node_modules": true},
		StaticFileExtensions: []string{".css", ".js", ".map", ".png", ".ico", ".jpg"},
		ReverseProxyRoute:    "/api",
	}
	reverseProxyURL := "http://localhost:3000"

	// Check that rootDir is a directory
	rootStat, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("root directory does not exists")
	}
	if !rootStat.IsDir() {
		return nil, fmt.Errorf("root directory is not a directory not exists")
	}

	// Read config file
	configContent, err := ioutil.ReadFile(filepath.Join(rootDir, "spago.json"))
	if err == nil {
		// Deserialize config
		var config Config
		err = json.Unmarshal(configContent, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize config: %v", err)
		}

		// Set if defined
		if config.Port != 0 {
			serv.Port = config.Port
		}
		if config.ExcludedPaths != nil {
			excludedPaths := make(map[string]bool)
			for _, p := range config.ExcludedPaths {
				excludedPaths[p] = true
			}
			serv.ExcludedPaths = excludedPaths
		}
		if config.StaticFileExtensions != nil {
			serv.StaticFileExtensions = config.StaticFileExtensions
		}
		if config.ReverseProxyRoute != "" {
			serv.ReverseProxyRoute = config.ReverseProxyRoute
		}
		if config.ReverseProxyURL != "" {
			reverseProxyURL = config.ReverseProxyURL
		}
	}

	// Create reverse proxy
	u, err := url.Parse(reverseProxyURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse ReverseProxyURL: %v", err)
	}
	serv.ReverseProxy = httputil.NewSingleHostReverseProxy(u)

	return serv, nil
}

func (serv *server) listenAndServe() {
	fmt.Printf("listening on http://localhost:%v\n", serv.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", serv.Port), serv))
}

func (serv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if serv.ReverseProxyRoute != "" && strings.HasPrefix(r.URL.Path, serv.ReverseProxyRoute) {
		serv.ReverseProxy.ServeHTTP(w, r)
		return
	}

	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	// Serve static files
	for _, extension := range serv.StaticFileExtensions {
		if strings.HasSuffix(r.URL.Path, extension) {
			http.ServeFile(w, r, filepath.Join(serv.RootDir, r.URL.Path))
			return
		}
	}

	var linkElements = make([]string, 0)
	var scriptElements = make([]string, 0)
	{
		err := unorderedWalk(serv.RootDir, serv.ExcludedPaths, "", func(path string) {
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

	htmlPath := filepath.Join(serv.RootDir, "index.html")

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
