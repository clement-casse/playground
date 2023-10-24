package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", "127.0.0.1", "The address used by the programm to expose API")
	port = flag.String("port", "8080", "The port the program will be listening to")
)

//go:embed static
var staticFS embed.FS

//go:embed templates
var templateFS embed.FS

//go:generate npm run build
var indexTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.ParseFS(templateFS, "templates/index.html"))
}

// indexTmplVars gather the variables that can be used in the index.html template
type indexTmplVars struct {
	Name string
}

func main() {
	// Get the values of Command Line Arguments (addr & port)
	flag.Parse()

	indexVars := &indexTmplVars{
		Name: "Unknown People",
	}

	// Define HTTP routes of the Application:
	// -> Everything under the URI Path `/static/` is a direct mapping of the ./static directory
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))
	// -> matches the static page "index" and Execute the HTML template before sending the response
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		indexTemplate.Execute(w, indexVars)
	})
	// -> Default behaviour, "/" matches everything, therefore reroute the / to index and raise 404 elsewhere
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/index", http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	})

	apiAddr := fmt.Sprintf("%s:%s", *addr, *port)
	log.Printf("Application Started and Listening on %s", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil)) // http.ListenAndServe is a blocking function, the main thread remains hanging there while serving HTTP requests
}
