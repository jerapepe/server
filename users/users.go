package users

import (
	"Project/auth"
	"Project/db"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	LastName     string `db:"last_name"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Role         string `db:"role"`
	ProfileImage []byte `db:"profile_image"`
}

//var connStr = "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"

func GetUser(username string) (User, bool, error) {
	us := User{}
	db := db.GetDB()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Queryx("SELECT id, name, last_name, username, email, role, profile_image FROM users where username = $1", username)
	if err != nil {
		return User{}, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var last_name string
		var username string
		var email string
		var role string
		var profile_image []byte

		err := rows.Scan(&id, &name, &last_name, &username, &email, &role, &profile_image)
		if err != nil {
			log.Fatal(err)
		}
		us = User{
			ID:           id,
			Name:         name,
			LastName:     last_name,
			Email:        email,
			Username:     username,
			Role:         role,
			ProfileImage: profile_image,
		}
	}
	return us, true, nil
}

func CreateUser(name string, last_name string, email string, username string, password string) (*User, bool, string) {
	if name == "" && last_name == "" && email == "" && username == "" && password == "" {
		fmt.Println("Null")
	} else {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			fmt.Println("Error al hashear la contraseña:", err)
			return nil, false, ""
		}
		UserN, err := userDB(name, last_name, email, username, hashedPassword)
		if err != nil {
			fmt.Println("Error", err)
			return nil, false, ""
		}
		token := auth.Cod(UserN.Username)
		return UserN, true, token
	}
	return nil, false, ""
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

func userDB(name string, last_name string, email string, username string, password []byte) (*User, error) {
	db := db.GetDB()
	_, err := db.Exec(`INSERT INTO users (name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5)`,
		name, last_name, username, email, password)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:     name,
		Email:    email,
		Username: username,
	}, nil
}

func getUserFromDB(username string) (*User, error) {
	db := db.GetDB()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Queryx("SELECT id, username, password FROM users where username = $1", username)
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

func DeleteUserFromDB(id int) error {
	db := db.GetDB()
	err := db.Ping()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func InsertSampleDat(name string, last_name string, username string, email string, password string) {
	db := db.GetDB()
	_, err := db.Exec(`
	    INSERT INTO users (name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5)
	`, name, last_name, username, email, password)
	if err != nil {
		log.Fatal(err)
	}
}

func AlterTable() {
	db := db.GetDB()
	_, err := db.Exec(`
		ALTER TABLE users
		ADD CONSTRAINT unique_username UNIQUE (username);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsersData() ([]byte, error) {
	db := db.GetDB()
	var users []User
	err := db.Select(&users, "SELECT name, last_name, username, email, role FROM users")
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func ConvertByteToJSON(byteData []byte) ([]User, error) {
	var users []User
	err := json.Unmarshal(byteData, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
