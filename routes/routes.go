package routes

import (
	"Project/auth"
	"Project/products"
	"Project/users"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const (
	endpointHome     = "/"
	endpointSignUp   = "/signup"
	endpointSignIn   = "/signin"
	endpointUser     = "/user"
	endpointAdmin    = "/admin"
	endpointProducts = "/products"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/user.html", "templates/login.html"))

type User struct {
	Username string
	Password string
	Token    string
}

var loggedInUser *User

func SetRoutes(router *mux.Router) {
	staticDir := "/static/files"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc(endpointHome, homeHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointSignUp, signUpHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointUser, userHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointAdmin, adminHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointProducts, productsHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointSignIn, signInHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)
	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
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

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
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

func signInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)
	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var formData struct {
				Username string
				Password string
				Token    string
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
				valid := auth.Decode(token)
				fmt.Println(valid)
				resultadoCh <- response
				return
			}
		}()

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
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

func userHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)
	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
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

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
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

func adminHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var formData struct {
				Username string
				Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}
			us, _, _ := users.GetUser(formData.Username)
			if us.Role == "admin" {
				userss, err := users.GetUsersData()
				if err != nil {
					fmt.Println(err)
				}
				usd, err := users.ConvertByteToJSON(userss)
				if err != nil {
					fmt.Println(err)
				}
				response := map[string]interface{}{"users": usd}
				resultadoCh <- response
				return
			}
		}()

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
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
		w.Write([]byte(`{"message":"dont send data"}`))
	}
}

func productsHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var formData struct {
				Username string
				Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}
			products, err := products.GetProducts()
			if err != nil {
				fmt.Println(err)
			}
			response := map[string]interface{}{"products": products}
			resultadoCh <- response
		}()

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
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
		w.Write([]byte(`{"message":"dont send data"}`))
	}
}
