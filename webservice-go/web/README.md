# WebApp

The `web` package contains a frontend application that can be embeded into the Go application, however, in theory, it can also be decoupled from the service binary and be served by a CDN.

## About the frontend application

The code of the frontend application can be found under the `app/` directory: it is built with Typescript, React and TailwindCSS; it also uses Vite.js for quick development iterations.
As of now the code does not do anything, it is just a frontend integration withing a Go application.

The frontend application is served through an embed file system by the go application, use the command `go generate -v ./...` in this directory to create a distribuable version of the application.


## References

- [Xe Iaso Blog post on using Tailwind CSS in Go][1]

[1]: https://xeiaso.net/blog/using-tailwind-go/
