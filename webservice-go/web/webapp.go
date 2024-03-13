package web

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
)

var frontendAppDir = "app/dist"

//go:generate npm --prefix ./app run build
//go:embed app/dist/*
var appdistFS embed.FS

// Handler returns an http.Handler that serves the web application
func Handler() http.Handler {
	webappFS, err := fs.Sub(appdistFS, frontendAppDir)
	if err != nil {
		panic("cannot create Sub File System for assets, has `go generate ./...` been run ?")
	}
	_ = mime.AddExtensionType(".js", "application/javascript")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webappFS)))
	return mux
}
