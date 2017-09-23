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

## Install

Download the release for your platform from the [release page](https://github.com/mushishi78/spago/releases).
Put in a folder in your PATH and rename to `spago`.

## Flags

```
spago -h
```

The `-h` flag will list the flag options available.

```
spago -PORT=3000
```

The `PORT` flag is used to set which port the dev server will listen on. The default is 8080.

## Build

To build from source, you will need to [install go](https://golang.org/doc/install).
Then use the `get` tool to download source in to your GOPATH.

```
go get -u github.com/mushishi78/spago
```

This will also build and install it into you go `bin` directory. To update with
changes made to the source, navigate to the project folder in your GOPATH and use
`go install`.
