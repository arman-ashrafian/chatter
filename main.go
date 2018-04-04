package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	jqueryCDN      = "https://code.jquery.com/jquery-3.3.1.min.js"
	materialcssCDN = "https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css"
	materialjsCDN  = "https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/js/materialize.min.js"
)

func main() {
	r := mux.NewRouter()

	// serve static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// handlers
	r.HandleFunc("/", indexHandler)

	// start server
	log.Println("Starting Server - Port 80")
	http.ListenAndServe(":80", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	fmt.Fprintf(w, indexHTML)
}

const indexHTML = `
<!DOCTYPE HTML>
<html>
  <head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="` + materialcssCDN + `">
	<link rel="stylesheet" href="/static/style.css">
    <title>Arman Ashrafian</title>
  </head>
  <body>
	<div id='root'>
		<h1 style="text-align: center"> Arman Ashrafian </h1>
	</div>
	<script src="` + jqueryCDN + `"></script>
	<script src="` + materialjsCDN + `"></script>
    <script src="/static/main.js"></script>
  </body>
</html>
`
