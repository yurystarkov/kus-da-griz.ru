package main

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/binary"
    "html/template"
    "net/http"
    "net/smtp"
    "os"
    "time"
)

var tmpl = template.Must(template.ParseFiles(
    "./templates/index.html",
    "./templates/contact_form.html",
    "./templates/navbar.html",
))

func main() {
    http.HandleFunc("/", index)

    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        tmpl.Execute(w, nil)
        return
    }

    name := r.FormValue("name")
    phone := r.FormValue("phone")

    sendMail([]byte(name + " " + phone))
    tmpl.Execute(w, struct{ Success bool }{true})
}

func sendMail(message []byte) {
    from := os.Getenv("MAIL_FROM")
    password := os.Getenv("MAIL_PASS")
    to := []string{os.Getenv("MAIL_TO")}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)
    smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}

func generateTOTP(secret string) (uint32, error) {
    // Convert the secret to byte array
    secretBytes := []byte(secret)

    // Get the number of 30-second intervals since the Unix epoch
    interval := time.Now().Unix() / 30

    // Convert the interval to byte array
    intervalBytes := make([]byte, 8)
    binary.BigEndian.PutUint64(intervalBytes, uint64(interval))

    // Create a new HMAC hasher
    hasher := hmac.New(sha1.New, secretBytes)

    // Write the interval to the hasher
    _, err := hasher.Write(intervalBytes)
    if err != nil {
        return 0, err
    }

    // Get the result of the HMAC and apply dynamic truncation
    hmacResult := hasher.Sum(nil)
    offset := hmacResult[len(hmacResult)-1] & 0x0F
    truncatedHash := hmacResult[offset : offset+4]

    // Convert the truncated hash to an integer
    code := binary.BigEndian.Uint32(truncatedHash) & 0x7FFFFFFF
    code = code % 1000000

    return code, nil
}
