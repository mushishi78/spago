package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

type failable interface {
	Name() string
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func tTempDir(t failable) (string, func()) {
	prefix := strings.Replace("spago-"+t.Name()+"-", "/", "-", -1)
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("failed to create a temp directory\n%v", err)
	}

	close := func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("failed to delete temp directory %v after test\n%v", dir, err)
		}
	}

	return dir, close
}

func tServerCreate(t failable, rootDir string) *server {
	serv, err := serverCreate(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	return serv
}

func tMkdir(t failable, dir string) {
	err := os.Mkdir(dir, 0600)
	if err != nil {
		t.Fatalf("failed to create dir\n%v", err)
	}
}

func tAddFile(t failable, filename string, content string) {
	err := ioutil.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		t.Fatalf("failed to add file\n%v", err)
	}
}

func tAddConfigFile(t failable, rootDir string, config Config) {
	content, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to serialize config file\n%v", err)
	}
	tAddFile(t, filepath.Join(rootDir, "spago.json"), string(content))
}

func tGetRequestEql(t failable, handler http.Handler, url string, status int, body string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create http request\n%v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != status {
		t.Errorf("expected status %v, got %v", status, rr.Code)
	}
	if rr.Body.String() != body {
		t.Errorf("\nexpected body\n\n%v\n\ngot\n\n%v", body, rr.Body.String())
		t.Errorf("\nlen: %v %v", len(body), len(rr.Body.String()))
	}
}
