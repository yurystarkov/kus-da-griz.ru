// simple CRUD with hierarchical
// filesystem
package data

import (
	"io/ioutil"
	"os"
)

type Product struct {
	ID          string
	Name        string
	Description string
	ImagePath   string
	Price       string
}

func ReadProduct(id string) Product {
	name       , _ := os.ReadFile("db/" + id + "/name")
	description, _ := os.ReadFile("db/" + id + "/description")
	imagePath  , _ := os.ReadFile("db/" + id + "/imagepath")
	price      , _ := os.ReadFile("db/" + id + "/price")

	return Product{
		id,
		string(name),
		string(description),
		string(imagePath),
		string(price),
	}
}

func ReadProducts() []Product {
	var products []Product

	productIDs, _ := ioutil.ReadDir("db/")

	for _, productID := range productIDs {
		products = append(products, ReadProduct(productID.Name()))
	}

	return products
}

func DeleteProduct(id string ) {
	os.RemoveAll("db/" + id)
}
