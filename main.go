package main

import (
	"net/http"

	"github.com/yurystarkov/kus-da-griz.ru/route"
)

func main() {
	http.HandleFunc("/"       , route.Index)
	http.HandleFunc("/catalog", route.Catalog)
	http.HandleFunc("/login"  , route.Login)
	http.HandleFunc("/logout" , route.Logout)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":8080", nil)
}
