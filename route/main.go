package route

import (
	"html/template"
	"net/http"
	"os"

	"github.com/pquerna/otp/totp"
	"github.com/yurystarkov/kus-da-griz.ru/data"
	"github.com/yurystarkov/kus-da-griz.ru/mail"
)

type CustomerInfo struct {
	Name  string
	Phone string
}

func Admin(w http.ResponseWriter, r *http.Request) {
	loginTmpl := template.Must(template.ParseFiles("./templates/admin.html"))
	if r.Method != http.MethodPost {
		loginTmpl.Execute(w, nil)
		return
	}
	if totp.Validate(r.FormValue("otp"), os.Getenv("SECRET")) {
		loginTmpl.Execute(w, struct{ Success bool }{true})
	}
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	catalogTmpl := template.Must(template.ParseFiles(
		"./templates/catalog.html",
		"./templates/navbar.html",
	))

	catalogTmpl.Execute(w, data.Products())
}

func Index(w http.ResponseWriter, r *http.Request) {
	indexTmpl := template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/contact_form.html",
		"./templates/navbar.html",
	))

	if r.Method != http.MethodPost {
		indexTmpl.Execute(w, nil)
		return
	}

	customer_info := CustomerInfo{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
	}

	mail.SendMailtoAdmin([]byte(customer_info.Name + " " + customer_info.Phone))
	indexTmpl.Execute(w, struct{ Success bool }{true})
}
