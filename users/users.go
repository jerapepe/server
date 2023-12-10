package users

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Name     string
	LastName string
	Email    string
	Username string
	Password string
}

func Users() User {
	return User{
		Name:     "Juan",
		LastName: "Perez",
		Email:    "juan@gmail.com",
		Username: "elpatitojuan",
		Password: "elpatitojuan",
	}
}

func GetUser(id int) User {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	us := User{}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, name, last_name, username, email, password FROM users where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var last_name string
		var username string
		var email string
		var password string

		err := rows.Scan(&id, &name, &last_name, &username, &email, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, last_name, username, email, password)
		us = User{
			ID:       id,
			Name:     name,
			LastName: last_name,
			Email:    email,
			Username: username,
			Password: password,
		}
	}
	return us
}
