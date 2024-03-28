package routes

import (
	"Project/auth"
	"Project/products"
	"Project/users"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
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
	endpointAccess   = "/access"
	endpointSearch   = "/search"
	endpointAdd      = "/add/product"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/user.html", "templates/login.html"))

type User struct {
	Username string
	Password string
	Token    string
}

func SetRoutes(router *mux.Router) {
	staticDir := "/static/files"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc(endpointHome, homeHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointSignUp, signUpHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointUser, userHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointAdmin, adminHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointAdd, addProductsHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointProducts, productsHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointAccess, accessHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointSearch, searchHadler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
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
			us, logged, token := users.CreateUser(formData.Name, formData.LastName, formData.Email, formData.Username, formData.Password)
			if logged {
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
			logged, us, token, _ := users.AuthenticationUser(formData.Username, formData.Password)
			if logged {
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
	if r.Method == "GET" {
		w.Write([]byte(`{"message": "Hello world"}`))
		err := templates.ExecuteTemplate(w, "login.html", "")
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
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type, authorization")

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	var wg sync.WaitGroup

	if r.Method == "POST" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token := r.Header.Get("Authorization")
			user, err := users.DecodeToken(token)
			if err != nil {
				response := map[string]interface{}{"user": "tokenExpired"}
				resultadoCh <- response
			}
			if user.Role == "admin" {
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
				Name string
			}
			if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
				errCh <- err
				return
			}
			products, err := products.GetProduct(formData.Name)
			if err != nil {
				fmt.Println(err)
			}
			response := map[string]interface{}{"product": products}
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
		wg.Add(1)
		go func() {
			defer wg.Done()
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
}

func accessHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type, authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	var wg sync.WaitGroup

	if r.Method == "POST" {
		token := r.Header.Get("Authorization")
		wg.Add(1)
		go func() {
			defer wg.Done()
			user, err := users.DecodeToken(token)
			if err != nil {
				response := map[string]interface{}{"user": "tokenExpired"}
				resultadoCh <- response
			}
			response := map[string]interface{}{"user": user}
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

func searchHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)

	var wg sync.WaitGroup

	if r.Method == "GET" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			query := r.URL.Query().Get("query")
			if query != "" {
				products, err := products.GetProduct(query)
				if err != nil {
					fmt.Println(err)
				}
				if products != nil {
					response := map[string]interface{}{"products": products}
					resultadoCh <- response
				}
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

func addProductsHadler(w http.ResponseWriter, r *http.Request) {
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

			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "Error al analizar el formulario", http.StatusBadRequest)
				errCh <- err
				return
			}

			name := r.FormValue("name")
			priceStr := r.FormValue("price")
			description := r.FormValue("description")
			sellerIDStr := r.FormValue("sellerId")

			file, handler, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "Error al acceder al archivo", http.StatusBadRequest)
				errCh <- err
				return
			}
			defer file.Close()
			filePath := fmt.Sprintf("uploads/%s", handler.Filename)
			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
				errCh <- err
				return
			}
			defer f.Close()
			io.Copy(f, file)

			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				http.Error(w, "Error al convertir el precio", http.StatusBadRequest)
				errCh <- err
				return
			}

			sellerID, err := strconv.Atoi(sellerIDStr)
			if err != nil {
				errCh <- err
				http.Error(w, "Error al convertir el ID del vendedor", http.StatusBadRequest)
				return
			}
			data := products.FormDatas{
				Name:         name,
				Price:        price,
				Description:  description,
				IDVendor:     sellerID,
				ProfileImage: []byte(filePath),
			}
			errs := products.AddProduct(products.FormDatas(data))
			if errs != nil {
				errCh <- err
			}
			response := map[string]interface{}{"products": "añadidos"}
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
}
