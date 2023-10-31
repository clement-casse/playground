# Go Single Page Application

## A Go backend App interacting with a Typescript Frontend

In this directory, I present both a *Go-powered backend application* and a *React frontend application* that are communicating together.
The React frontend application is bundled and made part of the Go application, so that it is served by the Go web-server.
The Go Backend embeds the minified code of the React frontend application and serve it alongside with a REST API.

### Repository Structure

#### Meta Files

These files are not required for the project per say, they allow me to quickly define and jump in particular development environments depending on the usecase.

- `flake.{nix,lock}` Definition of the development environement in a Nix Flake, in the same fashion as the other Directories in this repository.
- `.envrc` file used by [direnv](https://direnv.net/) to automatically load the nix flake when entering in this directory.

#### React Frontend Application

This application is made with the following tools (this impacts the structure of the application):

| Tool           | Purpose             |
| -------------- | ------------------- |
| *React*        | To dynamically generate WebPage content with JS |
| *Vite.js*      | To package and bundle our application while still providing a Hot-Reload mode for development. |
| *Tailwind.css* | To provide a nice look for the app without writting a single line of CSS .|

- `package.json` defines the manifest of the fontend application
- `{postcss,tailwind,vite}.config.js` configuration files of the multiples tools used to generate the application bundle.
- `webapp/` root of the React + Typescript frontend application source code.
- `dist/` root of the production optimised application obtained by running the `npm run build` or the `go generate ./...` command.
  This directory is used as embedded file system in the Go application.

#### Go Backend Server

In this playground, I try to apply [the *standard project layout* of a Go application][1]:

- `main.go` is the main entrypoint of the server
- `go.{mod,sum}` files defining this directory as a Go module alongs with its dependencies.
- `api/` a Go package defining the REST API of the server, look [here](https://github.com/golang-standards/project-layout/blob/master/api/README.md) for more details.
- `index.html` is both a Go Template and file used by Vite.js to generate the Web Application Code.

## Why the `go-spa` playground ?

## References

- [Xe Iaso Blog post on using Tailwind CSS in Go][2]
- [Standard Go Project Layout][1]

[1]: https://github.com/golang-standards/project-layout
[2]: https://xeiaso.net/blog/using-tailwind-go/
