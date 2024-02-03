package products

import (
	"Project/db"
	"fmt"

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
	var product []Product
	var err error
	db := db.GetDB()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx("SELECT * FROM product WHERE LOWER(name) = LOWER($1)", name)
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
		product = append(product, p)
	}
	return product, nil
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

type FormDatas struct {
	Name         string
	Price        float64
	Description  string
	IDVendor     int
	ProfileImage []byte
}

func AddProduct(data FormDatas) error {
	var err error

	db := db.GetDB()

	_, err = db.Exec(`
        INSERT INTO product (name, price, description, id_vendor, profile_image) VALUES ($1, $2, $3, $4, $5)
	`, data.Name, data.Price, data.Description, data.IDVendor, data.ProfileImage)
	if err != nil {
		return err
	}
	fmt.Println("Se agrego")
	return nil
}
