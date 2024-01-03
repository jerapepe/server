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
	"github.com/graphql-go/graphql"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/user.html", "templates/login.html"))

type User struct {
	Username string
	Password string
	Token    string
}

var loggedInUser *User

func main() {

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.Int},
			"name": &graphql.Field{Type: graphql.String},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := map[string]interface{}{
						"id":   1,
						"name": "Jera",
					}
					return user, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	staticDir := "/static/files"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc("/", HomeHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/signup", SignUpHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user", UserHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/signin", SignInHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	if r.Method == "POST" {
		go func() {
			var formData struct {
				Name     string
				LastName string
				Email    string
				Username string
				Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}
			fmt.Println(formData.Username)
			us, logged := users.CreateUser(formData.Name, formData.LastName, formData.Email, formData.Username, formData.Password)
			if logged {
				token := auth.Cod(us.Username)
				loggedInUser = &User{
					Username: us.Username,
					Password: us.Password,
					Token:    token,
				}
				response := map[string]interface{}{"token": token, "isVerified": true, "username": us.Username}
				resultadoCh <- response
				return
			}
		}()
		select {
		case res := <-resultadoCh:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusBadRequest)

		}
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	if r.Method == "POST" {
		go func() {
			var formData struct {
				Username string
				Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}
			logged, us, token, _ := users.Login(formData.Username, formData.Password)
			if logged {
				loggedInUser = &User{
					Username: us.Username,
					Password: us.Password,
					Token:    token,
				}
				response := map[string]interface{}{"token": token, "isVerified": true, "username": formData.Username}
				resultadoCh <- response
				return
			}
		}()
		select {
		case res := <-resultadoCh:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	if r.Method == "GET" {
		w.Write([]byte(`{"message": "Hello world"}`))
		err := templates.ExecuteTemplate(w, "login.html", loggedInUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func UserHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	if r.Method == "POST" {
		go func() {
			var formData struct {
				Username string
				Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}

			us, logged, _ := users.GetUser(formData.Username)
			if logged {
				token := auth.Cod(us.Username)
				response := map[string]interface{}{"token": token, "isVerified": true, "username": us.Username, "name": us.Name, "lastname": us.LastName, "email": us.Email}
				resultadoCh <- response
				return
			}
		}()

		select {
		case res := <-resultadoCh:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	if r.Method == "GET" {
		w.Write([]byte(`{"message":"hola"}`))
	}
}
