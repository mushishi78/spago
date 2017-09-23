# spago

Simple convention based tool for developing a single page application.

Running `spago` in a folder will set up an http server that serves all static
files in the folder, falling back to the provided `index.html` file for all
routes not found. In addition, the `index.html` file will have `link` and
`script` elements added for `css` and `js` files in the folder.

Typically development builds bundle assets just like in production. The main
benefit of serving the files directly is to work better with the development
tooling. For instance, [persistance workflows](https://developers.google.com/web/tools/setup/setup-workflow)
that allow changes made in the browser to be persisted directly to disk.

It's also fairly fast and responsive to changes.

