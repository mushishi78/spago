package main

import (
	"path/filepath"
	"testing"
)

func Test_serves_index_html_with_deps(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "user"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
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

func Test_fails_without_index_html(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "user"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tGetRequestEql(t, serv, "/cart", 500, "failed to read index.html\n")
}

func Test_fails_without_head_element(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "user"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
<html>
  <body>
    <div class="app"></div>
  </body>
</html>
`)
	tGetRequestEql(t, serv, "/cart", 500, "index.html does not have <head> element\n")
}

func Test_fails_without_body_element(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "user"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "user", "user.css"), ".user { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "user", "user.js"), "window.user = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
  </head>
</html>
`)
	tGetRequestEql(t, serv, "/cart", 500, "index.html does not have <body> element\n")
}

func Test_serves_static_assets(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tGetRequestEql(t, serv, "/cart/cart.css", 200, ".cart { border: 1px solid #666; }")
	tGetRequestEql(t, serv, "/cart/cart.js", 200, "window.cart = { total: 0 };")
}

func Test_serves_static_with_configured_extensions(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	tAddConfigFile(t, rootDir, Config{
		StaticFileExtensions: []string{".kik", ".jazz"},
	})
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.kik"), "kkkiikikikkkikiiiiik")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.jazz"), "doWa be ba dobob")
	tGetRequestEql(t, serv, "/cart/cart.kik", 200, "kkkiikikikkkikiiiiik")
	tGetRequestEql(t, serv, "/cart/cart.jazz", 200, "doWa be ba dobob")
}

func Test_excluded_files_not_added_to_index_html(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "node_modules"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "node_modules", "vendor.css"), ".vendor { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "node_modules", "vendor.js"), "window.vendor = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
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
  </head>
  <body>
    <div class="app"></div>
    <script src="/cart/cart.js"></script>
  </body>
</html>
`)
}

func Test_excluded_files_can_be_configured(t *testing.T) {
	rootDir, close := tTempDir(t)
	defer close()
	tAddConfigFile(t, rootDir, Config{
		ExcludedPaths: []string{"cake", "ingredients.js"},
	})
	serv := tServerCreate(t, rootDir)
	tMkdir(t, filepath.Join(rootDir, "cart"))
	tMkdir(t, filepath.Join(rootDir, "cake"))
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.css"), ".cart { border: 1px solid #666; }")
	tAddFile(t, filepath.Join(rootDir, "cart", "cart.js"), "window.cart = { total: 0 };")
	tAddFile(t, filepath.Join(rootDir, "cake", "cake.css"), ".cake { color: red; }")
	tAddFile(t, filepath.Join(rootDir, "cake", "cake.js"), "window.cake = { id: 'AE829X81PPD6' };")
	tAddFile(t, filepath.Join(rootDir, "ingredients.js"), "throw 'No INGREDIENTS!!!';")
	tAddFile(t, filepath.Join(rootDir, "index.html"), `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
  </head>
  <body>
    <div class="app"></div>
  </body>
</html>
`)
	tGetRequestEql(t, serv, "/arbitrary/path", 200, `<!DOCTYPE html>
<html>
  <head>
    <title>Example Project</title>
    <link href="/cart/cart.css" rel="stylesheet" type="text/css">
  </head>
  <body>
    <div class="app"></div>
    <script src="/cart/cart.js"></script>
  </body>
</html>
`)
}
