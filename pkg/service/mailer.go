package service

import (
	"os"
	"net/smtp"
	"log"
	"errors"
)

func SendMessage(sender string, message string) (error) {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
 	to := sender

 	msg := "From: " + from + "\n" +
	 "To: " + to + "\n" +
	 "Subject: " + message
 	err := smtp.SendMail("smtp.gmail.com:587",
 	smtp.PlainAuth("", from, password, "smtp.gmail.com"),
 	from, []string{to}, []byte(msg))

 	if err != nil {
	 	log.Printf("smtp error: %s", err)
		return errors.New("the recipient address is not a valid")
	}
	return nil
}