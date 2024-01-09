package route

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pquerna/otp/totp"
	"github.com/yurystarkov/kus-da-griz.ru/cmd/data"
	"github.com/yurystarkov/kus-da-griz.ru/cmd/mail"
)

var (
	key   = []byte(os.Getenv("AUTH_KEY"))
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
		loginTmpl.Execute(w, struct{ Success bool }{true})
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

func AdminDelete(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kus-da-griz-cookie")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	data.DeleteProduct(mux.Vars(r)["id"])
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminCreate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kus-da-griz-cookie")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	createTmpl := template.Must(template.ParseFiles("./templates/admin_create.html"))

	if r.Method != http.MethodPost {
		createTmpl.Execute(w, data.ReadProducts())
		return
	}

	id := uuid.NewString()
	name := r.FormValue("name")
	description := r.FormValue("description")
	price := r.FormValue("price")

	imageFile, imageHeader, err := r.FormFile("image")
	defer imageFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	imagePath := filepath.Join("assets/images/", id+filepath.Ext(imageHeader.Filename))
	dst, err := os.Create(imagePath)
	if err != nil {
		fmt.Println(err)
	}
	defer dst.Close()

	io.Copy(dst, imageFile)

	data.CreateProduct(id, name, price, description, imagePath)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func Index(w http.ResponseWriter, r *http.Request) {
	indexTmpl := template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/contact_form.html",
		"./templates/navbar.html",
	))

	if r.Method != http.MethodPost {
		indexTmpl.Execute(w, data.ReadProducts())
		return
	}

	customer_info := CustomerInfo{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
	}

	mail.SendMailtoAdmin([]byte(customer_info.Name + " " + customer_info.Phone))
	indexTmpl.Execute(w, nil)
}
