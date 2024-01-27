package main

import (
	"Project/db"
	"Project/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db.InitDB(connStr)

	router := mux.NewRouter()
	routes.SetRoutes(router)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Server 8000")
	log.Fatal(srv.ListenAndServe())
}
