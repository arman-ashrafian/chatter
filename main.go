package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	log.Println("Starting Server - Port 8080")
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	fmt.Fprintf(w, "HELLO WOLRD")
}
