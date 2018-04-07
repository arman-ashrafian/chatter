package main

import (
	"encoding/json"
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
	loginPath      = "./templates/login.html"
)

var indexTemplate *template.Template
var logins map[string]string

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // dev port
	}

	r := mux.NewRouter()

	// populate logins
	logins = make(map[string]string)
	logins["arman"] = "ash"

	// cache index template
	indexTemplate = template.Must(template.ParseFiles(basePath, indexPath))

	// serve static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// handlers
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)

	// handle server kill
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Printf("\nKilling Server\n\n")
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

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " --> " + r.URL.String())
	if r.Method == "POST" {
		var lf loginForm
		err := json.NewDecoder(r.Body).Decode(&lf)
		loginResp := loginResponse{"ok"}
		if err != nil {
			log.Println("Could not decode json")
			loginResp.Status = "error"
		}

		valid := checkLogin(lf.Username, lf.Password)

		if !valid {
			loginResp.Status = "error"
		}

		// send status
		// "okay" or "error"
		sendJSON(w, loginResp)
		return
	}

	// GET login page
	t, _ := template.ParseFiles(basePath, loginPath)
	t.ExecuteTemplate(w, "base", "")
}

func checkLogin(uname, pword string) bool {
	valid := false
	if logins[uname] == pword {
		valid = true
	}
	return valid
}
