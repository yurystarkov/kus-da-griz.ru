package mailing

import (
	"net/smtp"
	"os"
)

func SendMailtoAdmin(message []byte) {
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASS")
	to := []string{os.Getenv("MAIL_TO")}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}
