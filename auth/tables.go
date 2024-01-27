package auth

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

type Vendor struct {
	ID           int    `db:"id"`
	IDUser       int    `db:"id_user"`
	IDState      int    `db:"id_state"`
	IDLocation   int    `db:"id_location"`
	ProfileImage []byte `db:"profile_image"`
}

type State struct {
	ID                int    `db:"id"`
	IDUser            int    `db:"id_user"`
	State             string `db:"state"`
	STATE_FOR_SELLING string `db:"state_for_selling"`
	STATE_FOR_BUYING  string `db:"state_for_buying"`
	LastConnection    string `db:"last_connection"`
}

type Location struct {
	ID        int     `db:"id"`
	Longitude float64 `db:"longitude"`
	Latitude  float64 `db:"latitude"`
}

type Product struct {
	ID           int     `db:"id"`
	Name         string  `db:"name"`
	Price        float64 `db:"price"`
	Description  string  `db:"description"`
	IDVendor     int     `db:"id_vendor"`
	ProfileImage []byte  `db:"profile_image"`
}

type ProductCategory struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Order struct {
	ID          int     `db:"id"`
	IDUserBuyer int     `db:"id_user_buyer"`
	IDProduct   int     `db:"id_product"`
	OrderDate   string  `db:"order_date"`
	Status      string  `db:"status"`
	Quantity    int     `db:"quantity"`
	Total       float64 `db:"total"`
}

type PaymentMethod struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type FinancialTransaction struct {
	ID              int     `db:"id"`
	IDUserFrom      int     `db:"id_user_from"`
	IDUserTo        int     `db:"id_user_to"`
	Amount          float64 `db:"amount"`
	TransactionDate string  `db:"transaction_date"`
	Description     string  `db:"description"`
}

func CreateTables() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			last_name VARCHAR(100),
			username VARCHAR(100),
			email VARCHAR(100),
			password VARCHAR(100),
			role VARCHAR(100)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS location (
			id SERIAL PRIMARY KEY,
			longitude FLOAT,
			latitude FLOAT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS vendor (
		id SERIAL PRIMARY KEY,
		id_user INTEGER REFERENCES users(id),
		id_state INTEGER REFERENCES state(id),
		id_location INTEGER REFERENCES location(id)
	)
	`)

	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS state (
			id SERIAL PRIMARY KEY,
			id_user INTEGER REFERENCES users(id),
			state_for_selling VARCHAR(50),
			state_for_buying VARCHAR(50),
			last_connection TIMESTAMP
		);
		
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS product (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            price FLOAT,
            description TEXT,
            id_vendor INTEGER REFERENCES vendor(id)
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS product_categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100)
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            id_user_buyer INTEGER REFERENCES users(id),
            id_product INTEGER REFERENCES product(id),
            order_date TIMESTAMP,
            status VARCHAR(50),
            quantity INTEGER,
            total FLOAT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS payment_methods (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            description TEXT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS financial_transactions (
            id SERIAL PRIMARY KEY,
            id_user_from INTEGER REFERENCES users(id),
            id_user_to INTEGER REFERENCES users(id),
            amount FLOAT,
            transaction_date TIMESTAMP,
            description TEXT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tablas creadas correctamente")
}

func CheckTablesExistence() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		tableNames = append(tableNames, tableName)
	}

	fmt.Println("Tablas presentes en la base de datos:")
	for _, tableName := range tableNames {
		fmt.Println(tableName)
	}
}

func InsertSampleData() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
        INSERT INTO state (id_user, state_for_selling, state_for_buying, last_connection) VALUES (9, 'active', 'inactive', NOW())
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO location (longitude, latitude) VALUES (-74.006, 40.7128)
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO vendor (id_user, id_state, id_location) VALUES (9, 1, 1)
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO product (name, price, description, id_vendor) VALUES ('Pastes', 15.99, 'Pastes de arroz', 1)
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO product_categories (name) VALUES ('Postres')
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO orders (id_user_buyer, id_product, order_date, status, quantity, total) VALUES (9, 1, NOW(), 'pending', 1, 15.99)
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO payment_methods (name, description) VALUES ('Sample Method', 'Description of sample payment method')
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        INSERT INTO financial_transactions (id_user_from, id_user_to, amount, transaction_date, description) VALUES (9, 5, 100.0, NOW(), 'Sample transaction')
    `)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Datos insertados correctamente")
}

func ShowInsertedData() {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var stateData []State
	err = db.Select(&stateData, "SELECT * FROM state WHERE id_user = 9")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Datos insertados en la tabla 'state':")
	for _, data := range stateData {
		fmt.Println(data)
	}

	var locationData []Location
	err = db.Select(&locationData, "SELECT * FROM location")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Datos insertados en la tabla 'location':")
	for _, data := range locationData {
		fmt.Println(data)
	}

	var vendor []Vendor
	err = db.Select(&vendor, "SELECT * FROM vendor")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DAtos insertados en la tabla vendor")
	for _, data := range vendor {
		fmt.Println(data)
	}

	var payment_methods []PaymentMethod
	err = db.Select(&payment_methods, "SELECT * FROM payment_methods")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DAtos insertados en la tabla payment_methods")
	for _, data := range payment_methods {
		fmt.Println(data)
	}

	var financial_transactions []FinancialTransaction
	err = db.Select(&financial_transactions, "SELECT * FROM financial_transactions")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DAtos insertados en la tabla financial")
	for _, data := range financial_transactions {
		fmt.Println(data)
	}

	fmt.Println("Consulta de datos insertados finalizada")
}
