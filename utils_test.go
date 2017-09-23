package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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

func tRemoveFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	if err != nil {
		t.Fatalf("failed to remove file\n%v", err)
	}
}

func tFileEql(t *testing.T, filename string, expected string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("expected file %v to exist", filename)
		return
	}

	actual, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read output file\n%v", err)
	}

	if !bytes.Equal(actual, []byte(expected)) {
		t.Errorf("expected file %v content %v to be equal to %v", filename, actual, expected)
	}
}

func tHash(str string) string {
	array := md5.Sum([]byte(str))
	return hex.EncodeToString(array[:])
}

func tServerCreate(t *testing.T, cwd string) *Server {
	serv, err := ServerCreate(cwd)
	if err != nil {
		t.Fatalf("failed to create server\n%v", err)
	}
	return serv
}

func tServerClose(t *testing.T, serv *Server) {
	err := ServerClose(serv)
	if err != nil {
		t.Fatalf("failed to close server\n%v", err)
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
