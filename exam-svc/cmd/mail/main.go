package main

import (
	"log"

	"os"

	"github.com/joho/godotenv"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"
	mail "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/mailjet"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &config.MailjetConfig{
		API:  os.Getenv("MAILJET_API_KEY"),
		KEY:  os.Getenv("MAILJET_SECRET_KEY"),
		From: "hardwarerump04@gmail.com",
		Name: "Mailjet Pilot",
	}

	client := mailjet.NewMailjetClient(cfg.API, cfg.KEY)
	mailer := mail.NewMailjetClient(client, cfg.From, cfg.Name)

	err = mailer.SendTemplateEmail(
		"hardwarerump04@gmail.com",
		"passenger 1",
		mail.QuestionVerifiedTemplate,
	)
	if err != nil {
		log.Fatal("Failed to send email:", err)
	}

	log.Println("Email sent successfully!")
}
