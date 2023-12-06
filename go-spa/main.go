package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/clement-casse/playground/go-spa/api"
)

var (
	Version   = "v0-dirty"
	StartTime = time.Now()
)

var (
	addr = flag.String("addr", "127.0.0.1", "The address used by the programm to expose API")
	port = flag.String("port", "8080", "The port the program will be listening to")
)

//go:generate npm run build
//go:embed dist
var distFS embed.FS

var indexTemplate = template.Must(template.ParseFS(distFS, "dist/index.html"))

// indexTmplVars gather the variables that can be used in the index.html template
type indexTmplVars struct {
	Title          string
	BackendVersion string
	BackendUptime  time.Duration
}

func main() {
	// Get the values of Command Line Arguments (addr & port)
	flag.Parse()

	indexVars := &indexTmplVars{
		Title:          "What a Marvelous title !",
		BackendVersion: Version,
	}

	assetsFS, err := fs.Sub(distFS, "dist/assets")
	if err != nil {
		log.Fatalf("cannot create Sub File System for asssets, has `go generate ./...` been run ?")
	}

	// Define HTTP router for the Application
	mux := http.NewServeMux()

	// And then map the routes:
	// -> Everything under the URI Path `/assets/` is a direct mapping of the ./static directory
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.FS(assetsFS))))
	// -> Everything matching the URI Path `/api/` is managed by the Router of the `api` module
	mux.Handle("/api/", http.StripPrefix("/api", api.NewAPIRouter()))
	// -> matches the static page "index" and Execute the HTML template before sending the response
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		indexVars.BackendUptime = time.Since(StartTime)
		err := indexTemplate.Execute(w, indexVars)
		if err != nil {
			log.Print("executing index template returned an error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	// -> Default behaviour, "/" matches everything, therefore reroute the / to index and raise 404 elsewhere
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/index", http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	})

	apiAddr := fmt.Sprintf("%s:%s", *addr, *port)
	log.Printf("Application Started and Listening on %s", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, mux)) // http.ListenAndServe is a blocking function, the main thread remains hanging there while serving HTTP requests
}
