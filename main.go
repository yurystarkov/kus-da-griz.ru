package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/yurystarkov/kus-da-griz.ru/mailing"

	"github.com/pquerna/otp/totp"

	bolt "go.etcd.io/bbolt"
)

type CustomerInfo struct {
	Name  string
	Phone string
}

type ProductData struct {
	Name       string
	Descrption string
	ImagePath  string
	Price      string
}

var indexTmpl = template.Must(template.ParseFiles(
	"./templates/index.html",
	"./templates/contact_form.html",
	"./templates/navbar.html",
))

var catalogTmpl = template.Must(template.ParseFiles(
	"./templates/catalog.html",
	"./templates/navbar.html",
))

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/catalog", catalog)
	http.HandleFunc("/admin", admin)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":8080", nil)
}

func admin(w http.ResponseWriter, r *http.Request) {
	loginTmpl := template.Must(template.ParseFiles("./templates/admin.html"))
	if r.Method != http.MethodPost {
		loginTmpl.Execute(w, nil)
		return
	}
	if totp.Validate(r.FormValue("otp"), os.Getenv("SECRET")) {
		loginTmpl.Execute(w, struct{ Success bool }{true})
	}
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("products.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Products"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	defer db.Close()
	return db, nil
}

func catalog(w http.ResponseWriter, r *http.Request) {
	catalogTmpl.Execute(w, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		indexTmpl.Execute(w, nil)
		return
	}

	customer_info := CustomerInfo{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
	}

	mailing.SendMailtoAdmin([]byte(customer_info.Name + " " + customer_info.Phone))
	indexTmpl.Execute(w, struct{ Success bool }{true})
}
