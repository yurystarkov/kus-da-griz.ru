package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ProductData struct {
	Name       string `json:"name"`
	Description string `json:"description"`
	ImagePath  string `json:"image_path"`
	Price      string `json:"price"`
}

func Products() []ProductData {
	var products []ProductData
	productFiles, err := ioutil.ReadDir("db")
	if err != nil {
		log.Println(err)
	}

	for _, productFile := range productFiles {
		productFileContent, err := ioutil.ReadFile("db/" + productFile.Name())
		if err != nil {
			log.Println(err)
		}
		var product ProductData
		err = json.Unmarshal([]byte(productFileContent), &product)
		if err != nil {
			log.Println(err)
		}
		products = append(products, product)
	}

	return products
}
