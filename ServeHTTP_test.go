package main

import (
	"path/filepath"
	"testing"
)

func TestServeHTTP_serves_index_html_with_deps(t *testing.T) {
	cwd, close := tTempDir(t)
	defer close()
	serv := &server{cwd}
	tMkdir(t, filepath.Join(cwd, "cart"))
	tMkdir(t, filepath.Join(cwd, "user"))
	tAddFile(t, filepath.Join(cwd, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(cwd, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(cwd, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(cwd, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(cwd, "index.html"), `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
  </head>
  <body>
    <div class="app"></div>
  </body>
</html>
`)
	tGetRequestEql(t, serv, "/cart", 200, `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
    <link href="/cart/cart.css" rel="stylesheet" type="text/css">
    <link href="/user/user.css" rel="stylesheet" type="text/css">
  </head>
  <body>
    <div class="app"></div>
    <script src="/cart/cart.js"></script>
    <script src="/user/user.js"></script>
  </body>
</html>
`)
}

func TestServeHTTP_fails_without_index_html(t *testing.T) {
	cwd, close := tTempDir(t)
	defer close()
	serv := &server{cwd}
	tMkdir(t, filepath.Join(cwd, "cart"))
	tMkdir(t, filepath.Join(cwd, "user"))
	tAddFile(t, filepath.Join(cwd, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(cwd, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(cwd, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(cwd, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tGetRequestEql(t, serv, "/cart", 500, "failed to read index.html\n")
}

func TestServeHTTP_fails_without_head_element(t *testing.T) {
	cwd, close := tTempDir(t)
	defer close()
	serv := &server{cwd}
	tMkdir(t, filepath.Join(cwd, "cart"))
	tMkdir(t, filepath.Join(cwd, "user"))
	tAddFile(t, filepath.Join(cwd, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(cwd, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(cwd, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(cwd, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(cwd, "index.html"), `<!DOCTYPE html>
<html>
  <body>
    <div class="app"></div>
  </body>
</html>
`)
	tGetRequestEql(t, serv, "/cart", 500, "index.html does not have <head> element\n")
}

func TestServeHTTP_fails_without_body_element(t *testing.T) {
	cwd, close := tTempDir(t)
	defer close()
	serv := &server{cwd}
	tMkdir(t, filepath.Join(cwd, "cart"))
	tMkdir(t, filepath.Join(cwd, "user"))
	tAddFile(t, filepath.Join(cwd, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(cwd, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(cwd, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(cwd, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(cwd, "index.html"), `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
  </head>
</html>
`)
	tGetRequestEql(t, serv, "/cart", 500, "index.html does not have <body> element\n")
}
