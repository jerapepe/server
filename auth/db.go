package auth

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DataBase() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, name, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var username string
		var password string
		err := rows.Scan(&id, &name, &username, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, username, password)
	}
}
