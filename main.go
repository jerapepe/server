package main

import (
	"Project/auth"
	"Project/users"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/user.html", "templates/login.html"))

type User struct {
	Username string
	Password string
	Token    string
}

var loggedInUser *User

func main() {
	router := mux.NewRouter()

	staticDir := "/static/files"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc("/", HomeHandler).Methods("GET", "POST")
	router.HandleFunc("/signup", SignUpHandler).Methods("POST", "GET")
	router.HandleFunc("/user", UserHadler).Methods("GET", "POST")
	router.HandleFunc("/signin", SignInHandler).Methods("GET", "POST")

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Server 8000")
	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		errs := r.ParseForm()
		if errs != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		name := r.Form.Get("name")
		lastName := r.Form.Get("lastname")
		email := r.Form.Get("email")
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		us := users.CreateUser(name, lastName, email, username, password)
		token := auth.Cod(us.Username)
		loggedInUser = &User{
			Username: us.Username,
			Password: us.Password,
			Token:    token,
		}
		err := templates.ExecuteTemplate(w, "user.html", loggedInUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "login.html", loggedInUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		errs := r.ParseForm()
		if errs != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		logged, us, token, _ := users.Login(username, password)
		if logged {
			loggedInUser = &User{
				Username: us.Username,
				Password: us.Password,
				Token:    token,
			}
			err := templates.ExecuteTemplate(w, "user.html", loggedInUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func UserHadler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	}
	if r.Method == "POST" {
		errs := r.ParseForm()
		if errs != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		err := templates.ExecuteTemplate(w, "user.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
