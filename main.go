package main

import (
    "html/template"
    "net/http"
    "fmt"
)

type CustomerInfo struct {
    Name  string
    Phone string
}

func main() {
    tmpl := template.Must(template.ParseFiles("./templates/index.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        details := CustomerInfo{
            Name:  r.FormValue("name"),
            Phone: r.FormValue("phone"),
        }

        fmt.Println(details.Name)
        fmt.Println(details.Phone)

        tmpl.Execute(w, struct{ Success bool }{true})
    })

    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))
    http.ListenAndServe(":8080", nil)
}
