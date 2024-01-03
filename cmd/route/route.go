package route

import (
	"html/template"
	"net/http"
	"os"

	"github.com/pquerna/otp/totp"
	"github.com/yurystarkov/kus-da-griz.ru/cmd/data"
	"github.com/yurystarkov/kus-da-griz.ru/cmd/mail"
	"github.com/gorilla/sessions"
)

var (
    key = []byte(os.Getenv("AUTH_KEY"))
    store = sessions.NewCookieStore(key)
)

type CustomerInfo struct {
	Name  string
	Phone string
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kus-da-griz-cookie")
	loginTmpl := template.Must(template.ParseFiles("./templates/login.html"))
	if r.Method != http.MethodPost {
		loginTmpl.Execute(w, nil)
		return
	}
	if totp.Validate(r.FormValue("otp"), os.Getenv("OTP_SECRET")) {
		session.Values["authenticated"] = true
		session.Save(r, w)
		loginTmpl.Execute(w, struct{Success bool}{true})
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "kus-da-griz-cookie")
    session.Values["authenticated"] = false
    session.Save(r, w)
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	catalogTmpl := template.Must(template.ParseFiles(
		"./templates/catalog.html",
		"./templates/navbar.html",
	))

	catalogTmpl.Execute(w, data.ReadProducts())
}

func Admin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kus-da-griz-cookie")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	catalogTmpl := template.Must(template.ParseFiles("./templates/admin.html"))
	catalogTmpl.Execute(w, data.ReadProducts())
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
