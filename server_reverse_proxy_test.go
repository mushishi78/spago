package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

type apiHandler struct{}

func (ah *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API server says hi: "+r.URL.Path)
}

func Test_uses_reverse_proxy(t *testing.T) {
	apiServ := &http.Server{Addr: ":3000", Handler: &apiHandler{}}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		apiServ.ListenAndServe()
		wg.Done()
	}()

	go func() {
		rootDir, close := tTempDir(t)
		defer close()
		serv := tServerCreate(t, rootDir)
		tGetRequestEql(t, serv, "/api/hello", 200, "API server says hi: /api/hello\n")
		apiServ.Shutdown(nil)
	}()

	wg.Wait()
}
func Test_uses_configured_reverse_proxy_url(t *testing.T) {
	apiServ := &http.Server{Addr: ":4040", Handler: &apiHandler{}}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		apiServ.ListenAndServe()
		wg.Done()
	}()

	go func() {
		rootDir, close := tTempDir(t)
		defer close()
		tAddConfigFile(t, rootDir, Config{
			ReverseProxyURL: "http://localhost:4040",
		})
		serv := tServerCreate(t, rootDir)
		tGetRequestEql(t, serv, "/api/hello", 200, "API server says hi: /api/hello\n")
		apiServ.Shutdown(nil)
	}()

	wg.Wait()
}
func Test_uses_configured_reverse_proxy_route(t *testing.T) {
	apiServ := &http.Server{Addr: ":3000", Handler: &apiHandler{}}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		apiServ.ListenAndServe()
		wg.Done()
	}()

	go func() {
		rootDir, close := tTempDir(t)
		defer close()
		tAddConfigFile(t, rootDir, Config{
			ReverseProxyRoute: "/backend",
		})
		serv := tServerCreate(t, rootDir)
		tGetRequestEql(t, serv, "/backend/hello", 200, "API server says hi: /backend/hello\n")
		apiServ.Shutdown(nil)
	}()

	wg.Wait()
}
