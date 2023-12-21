package main

import (
    "html/template"
    "net/http"
    "net/smtp"
    "os"
)

func main() {
    tmpl := template.Must(template.ParseFiles(
        "./templates/index.html",
        "./templates/contact_form.html",
        "./templates/navbar.html",
    ))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        name := r.FormValue("name")
        phone := r.FormValue("phone")

        sendMail([]byte(name + " " + phone))
        tmpl.Execute(w, struct{ Success bool }{true})
    })

    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))
    http.ListenAndServe(":8080", nil)
}

func sendMail(message []byte)  {
    from := os.Getenv("MAIL_FROM")
    password := os.Getenv("MAIL_PASS")
    to := []string{os.Getenv("MAIL_TO")}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)
    smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, to, message)
}
