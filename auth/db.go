package auth

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DataBase() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var users []User
	err = db.Select(&users, "SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Datos insertados en la tabla 'users':")
	for _, data := range users {
		fmt.Println(data)
	}
}

func ShowTableSchema(tableName string) {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", tableName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("Diseño de la tabla '%s':\n", tableName)
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Columna: %s - Tipo de dato: %s\n", columnName, dataType)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
