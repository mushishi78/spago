package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

/*
$ go test -bench=. -run=^a
BenchmarkServeHTTP_index_file/Small_Project-4               2000            988036 ns/op
BenchmarkServeHTTP_index_file/Medium_Project-4               300           3500050 ns/op
BenchmarkServeHTTP_index_file/Large_Project-4                200           7031795 ns/op
PASS
ok      github.com/mushishi78/spago     10.590s
*/

func Benchmark_server_index_file(b *testing.B) {
	for _, tt := range []struct {
		Name      string
		RuneRange string
	}{
		{"Small Project", "abc"},
		{"Medium Project", "abcde"},
		{"Large Project", "abcdefg"},
	} {
		b.Run(tt.Name, func(b *testing.B) {
			rootDir, close := tTempDir(b)
			defer close()

			// Create nested project
			for _, r1 := range tt.RuneRange {
				tMkdir(b, filepath.Join(rootDir, string(r1)))

				for _, r2 := range tt.RuneRange {
					tMkdir(b, filepath.Join(rootDir, string(r1), string(r2)))

					for _, r3 := range tt.RuneRange {
						tAddFile(b, filepath.Join(rootDir, string(r1), string(r2), string(r3)+".css"), "Lorem ipsum dolor sit amet")
						tAddFile(b, filepath.Join(rootDir, string(r1), string(r2), string(r3)+".js"), "Lorem ipsum dolor sit amet")
					}
				}
			}
			tAddFile(b, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
        <html>
        <head>
            <title>Example Project</title>
        </head>
        <body>
            <div class="app"></div>
        </body>
        </html>
        `)

			serv := tServerCreate(b, rootDir)
			req, err := http.NewRequest("GET", "/test/route", nil)
			if err != nil {
				b.Fatalf("failed to create http request\n%v", err)
			}
			// Benchmark the GET request
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				serv.ServeHTTP(httptest.NewRecorder(), req)
			}
		})
	}
}
