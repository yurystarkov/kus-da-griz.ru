package main

import (
    "html/template"
    "net/http"
    "net/smtp"
    "os"
    "fmt"
)

type CustomerInfo struct {
    Name  string
    Phone string
}

func main() {
    tmpl := template.Must(template.ParseFiles("./templates/index.html", "./templates/_navbar.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        details := CustomerInfo{
            Name:  r.FormValue("name"),
            Phone: r.FormValue("phone"),
        }

        sendMail([]byte(details.Name + " " + details.Phone))
        tmpl.Execute(w, struct{ Success bool }{true})
    })

    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))
    http.ListenAndServe(":8080", nil)
}

func sendMail(message []byte)  {
    from := os.Getenv("MAIL_FROM")
    password := os.Getenv("MAIL_PASS")

    to := []string{
        os.Getenv("MAIL_TO"),
    }

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, to, message)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Email Sent Successfully!")
}
