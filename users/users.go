package users

import (
	"Project/auth"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	LastName string
	Email    string
	Username string
	Password string
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

func CreateUser(name string, last_name string, email string, username string, password string) (*User, bool) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error al hashear la contraseña:", err)
		return nil, false
	}
	UserN, rr := UserDB(name, last_name, email, username, hashedPassword)
	if rr != nil {
		fmt.Println("Error", err)
		return nil, false
	}
	err = comparePasswords(hashedPassword, password)
	if err != nil {
		fmt.Println("Las contraseñas no coinciden")
		return nil, false
	}
	return UserN, true
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func comparePasswords(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err
}

func UserDB(name string, last_name string, email string, username string, password []byte) (*User, error) {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO users (name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5)`,
		name, last_name, username, email, password)
	if err != nil {
		log.Fatal(err)
	}
	return &User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: string(password),
	}, nil
}

func Login(username, password string) (bool, *User, string, error) {
	user, err := getUserFromDB(username)
	if err != nil {
		return false, nil, "", err
	}
	if user == nil {
		return false, nil, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil, "", err
	}
	token := auth.Cod(user.Username)
	return true, user, token, nil
}

func getUserFromDB(username string) (*User, error) {
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
	rows, err := db.Query("SELECT id, username, password FROM users where username = $1", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	userFromDB := &User{}
	for rows.Next() {
		var id int
		var username string
		var password string

		err := rows.Scan(&id, &username, &password)
		if err != nil {
			log.Fatal(err)
		}
		userFromDB = &User{
			ID:       id,
			Username: username,
			Password: password,
		}
	}
	return userFromDB, nil
}
