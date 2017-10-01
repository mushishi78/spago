package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

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
			cwd, close := tTempDir(b)
			defer close()

			// Create nested project
			for _, r1 := range tt.RuneRange {
				tMkdir(b, filepath.Join(cwd, string(r1)))

				for _, r2 := range tt.RuneRange {
					tMkdir(b, filepath.Join(cwd, string(r1), string(r2)))

					for _, r3 := range tt.RuneRange {
						tAddFile(b, filepath.Join(cwd, string(r1), string(r2), string(r3)+".css"), lorem)
						tAddFile(b, filepath.Join(cwd, string(r1), string(r2), string(r3)+".js"), lorem)
					}
				}
			}
			tAddFile(b, filepath.Join(cwd, "index.html"), `<!DOCTYPE html>
        <html>
        <head>
            <title>Example Project</title>
        </head>
        <body>
            <div class="app"></div>
        </body>
        </html>
        `)

			serv := tServerCreate(b, cwd, 3000)
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
