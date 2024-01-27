package products

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Product struct {
	ID           int     `db:"id"`
	Name         string  `db:"name"`
	Price        float64 `db:"price"`
	Description  string  `db:"description"`
	IDVendor     int     `db:"id_vendor"`
	ProfileImage []byte  `db:"profile_image"`
}

func GetProducts() ([]Product, error) {
	connStr := "host=192.168.0.73 port=5432 user=postgres dbname=marketupi password=mi_contraseña sslmode=disable"
	var products []Product
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx("SELECT id, name, price, description, id_vendor, profile_image FROM product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		fmt.Println(p.ID, p.Name, p.Price, p.Description, p.IDVendor, p.ProfileImage)
		products = append(products, p)
	}
	return products, nil
}
