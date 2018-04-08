package main

import (
	"digitaloceanssampleapp/chat"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	materialcssCDN = "https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css"
	vuejsCDN       = "https://cdn.jsdelivr.net/npm/vue"
	basePath       = "./templates/base.html"
	indexPath      = "./templates/index.html"
	loginPath      = "./templates/login.html"
)

var (
	indexTemplate *template.Template
	logins        map[string]string
	sessionkey    = []byte("secret-key")
	store         = sessions.NewCookieStore(sessionkey)
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // dev port
	}

	r := mux.NewRouter()

	// populate logins
	logins = make(map[string]string)
	logins["arman"] = "ash"
	logins["andrew"] = "ilovearman"
	logins["luis"] = "ilovearman"
	logins["westin"] = "ilovearman"

	// cache index template
	indexTemplate = template.Must(template.ParseFiles(basePath, indexPath))

	// serve static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// start chat server
	server := chat.NewServer(r)
	go server.Listen()

	// handlers
	r.HandleFunc("/", reqLogin(indexHandler))
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

func reqLogin(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")
		if auth, ok := session.Values["auth"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		f(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " --> " + r.URL.String())

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

		if valid {
			// set auth cookie
			session, _ := store.Get(r, "auth")
			session.Values["auth"] = true
			session.Save(r, w)
		} else {
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
	realpass, ok := logins[uname]
	if !ok {
		return false
	}
	return pword == realpass
}
