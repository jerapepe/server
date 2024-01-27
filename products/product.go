package products

import (
	"Project/db"

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

func GetProduct(name string) ([]Product, error) {
	var products []Product
	var err error

	db := db.GetDB()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx("SELECT id, name, price, description, id_vendor, profile_image FROM product where name = $1", name)
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
		products = append(products, p)
	}
	return products, nil
}

func GetProducts() ([]Product, error) {
	var products []Product

	var err error

	db := db.GetDB()

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
		products = append(products, p)
	}
	return products, nil
}
