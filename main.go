package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yurystarkov/kus-da-griz.ru/route"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/"       , route.Index)
	r.HandleFunc("/catalog", route.Catalog)
	r.HandleFunc("/login"  , route.Login)
	r.HandleFunc("/logout" , route.Logout)
	r.HandleFunc("/admin"  , route.Admin)

	http.Handle("/", r)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":8080", nil)
}
