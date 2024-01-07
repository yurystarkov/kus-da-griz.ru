package mail

import (
	"log"
	"net/smtp"
	"os"
	_ "github.com/joho/godotenv/autoload"
)

func SendMailtoAdmin(message []byte) {
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASS")
	to := []string{os.Getenv("MAIL_TO")}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println(err)
	}
}
