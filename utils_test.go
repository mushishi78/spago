package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func tTempDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "spago-"+t.Name()+"-")
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

func tServerCreate(t *testing.T, cwd string, apiPort int) *server {
	serv, err := serverCreate(cwd, apiPort)
	if err != nil {
		t.Fatal(err)
	}
	return serv
}

func tMkdir(t *testing.T, dir string) {
	err := os.Mkdir(dir, 0600)
	if err != nil {
		t.Fatalf("failed to create dir\n%v", err)
	}
}

func tAddFile(t *testing.T, filename string, content string) {
	err := ioutil.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		t.Fatalf("failed to add file\n%v", err)
	}
}

func tGetRequestEql(t *testing.T, handler http.Handler, url string, status int, body string) {
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
