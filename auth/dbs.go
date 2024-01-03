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
	fmt.Println("Base de datos creada")
}

func UpdateUserRoles() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE users
		SET role = 
			CASE 
				WHEN id = 5 THEN 'admin'
				ELSE 'usuario'
			END
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tabla alterada")
}
