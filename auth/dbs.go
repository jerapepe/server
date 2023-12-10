package auth

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func CreateDB() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketup password=mi_contraseña sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE marketUPI")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Base de datos creada correctamente")
}

func CreateTable() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		last_name VARCHAR(100),
		username VARCHAR(100),
		email VARCHAR(100),
		password VARCHAR(100)
	)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tabla creada correctamente")
}

func CreateR() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO users (name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5)`,
		"JeraPepe", "Ruiz", "Jera", "ruizpepe402@gmail.com", "Elpatitojuan")
	if err != nil {
		log.Fatal(err)
	}
}
