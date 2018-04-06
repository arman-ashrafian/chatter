package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

const (
	materialcssCDN = "https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css"
	vuejsCDN       = "https://cdn.jsdelivr.net/npm/vue"
	basePath       = "./templates/base.html"
	indexPath      = "./templates/index.html"
)

var indexTemplate *template.Template

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // dev port
	}

	r := mux.NewRouter()

	// cache index template
	indexTemplate = template.Must(template.ParseFiles(basePath, indexPath))

	// serve static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// handlers
	r.HandleFunc("/", indexHandler)

	// handle server kill
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nKilling Server\n")
		shutdownServer()
	}()

	// start server
	log.Println("Starting Server - Port " + port)
	http.ListenAndServe(":"+port, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	indexTemplate.ExecuteTemplate(w, "base", "Arman Ashrafian")
}

func shutdownServer() {
	fmt.Println("Server shutdown")
	os.Exit(0)
}
