# spago

Simple convention based tool for developing a single page application.

Running `spago` in a folder will set up an http server that serves all static
files in the folder, falling back to the provided `index.html` file for all
routes not found. In addition, the `index.html` file will have `link` and
`script` elements added for `css` and `js` files in the folder.

Typically development builds bundle assets just like in production. The main
benefit of serving the files directly is to use [persistance workflows](https://developers.google.com/web/tools/setup/setup-workflow)
that allow changes made in the browser to be persisted directly to disk.

It's also fairly fast and responsive to changes.

## API forwarding

By default, all http requests under the `/api` route will be forwarded to
`http://localhost:3000`. This allows a backend server to be run side by side
the `spago` dev server, whilst be considered the same domain as far as the
browser is concerned.

## Install

Download the release for your platform from the [release page](https://github.com/mushishi78/spago/releases).
Put in a folder in your PATH `spago`.

## Usage

```
spago [path]
```

If a path is not provided, the current working directory will be used.

## Config

Spago will look for a `spago.json` file in the directory that it is run in.
If no file if found a default configuration as follows will be used:

```
{
    "port": 8080,
    "excludedPaths": ["node_modules"],
    "staticFileExtensions": [".css", ".js", ".map", ".png", ".ico", ".jpg"],
    "reverseProxyUrl": "http://localhost:3000",
    "reverseProxyRoute": "/api"
}
```

## Build

To build from source, you will need to [install go](https://golang.org/doc/install).
Then use the `get` tool to download source in to your GOPATH.

```
go get -u github.com/mushishi78/spago
```

This will also build and install it into you go `bin` directory. To update with
changes made to the source, navigate to the project folder in your GOPATH and use
`go install`.
