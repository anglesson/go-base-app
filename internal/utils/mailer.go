package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

func SendPasswordResetEmail(email string, token string) error {
	log.Print("Sending email")
	m := mail.NewMsg()

	if err := m.From("toni.sender@example.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}

	if err := m.To("tina.recipient@example.com"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	m.Subject("Recover your password!")

	m.SetBodyString(mail.TypeTextHTML, "<a href='"+os.Getenv("CLIENT_URL")+"/reset-password?token="+token+"'>Click here</a>")

	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	userName := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	c, err := mail.NewClient(host, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(userName), mail.WithPassword(password))

	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

	return nil
}
