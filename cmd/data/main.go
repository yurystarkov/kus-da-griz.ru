package data

import (
	"io/ioutil"
	"os"
)

type ProductData struct {
	ID          string
	Name        string
	Description string
	ImagePath   string
	Price       string
}

func ReadProduct(id string) ProductData {
	name       , _ := os.ReadFile("db/" + id + "/name")
	description, _ := os.ReadFile("db/" + id + "/description")
	imagePath  , _ := os.ReadFile("db/" + id + "/imagepath")
	price      , _ := os.ReadFile("db/" + id + "/price")

	return ProductData{
		id,
		string(name),
		string(description),
		string(imagePath),
		string(price),
	}
}

func ReadProducts() []ProductData {
	var products []ProductData

	productIDs, _ := ioutil.ReadDir("db/")

	for _, productID := range productIDs {
		products = append(products, ReadProduct(productID.Name()))
	}

	return products
}
