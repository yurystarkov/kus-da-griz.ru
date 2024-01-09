// simple CRUD with hierarchical filesystem
// directory names are IDs, filenames are
// keys and their contents are values
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
	name, _ := os.ReadFile("db/" + id + "/name")
	description, _ := os.ReadFile("db/" + id + "/description")
	imagePath, _ := os.ReadFile("db/" + id + "/imagepath")
	price, _ := os.ReadFile("db/" + id + "/price")

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

func CreateProduct(
	id string,
	name string,
	price string,
	description string,
	imagePath string,
) {
	os.Mkdir("db/"+id, os.ModePerm)

	err := os.WriteFile("db/"+id+"/name", []byte(name), 0666)
	if err != nil {
		panic(err)
	}

	os.WriteFile("db/"+id+"/price", []byte(price), 0666)
	if err != nil {
		panic(err)
	}

	os.WriteFile("db/"+id+"/description", []byte(description), 0666)
	if err != nil {
		panic(err)
	}

	os.WriteFile("db/"+id+"/imagepath", []byte(imagePath), 0666)
	if err != nil {
		panic(err)
	}
}

func DeleteProduct(id string) {
	os.RemoveAll("db/" + id)
}
