package main

import (
	"log"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"
	mail "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/mailjet"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	mailjetCfg := cfg.Mailjet

	client := mailjet.NewMailjetClient(mailjetCfg.API, mailjetCfg.KEY)
	mailer := mail.NewMailjetClient(client, mailjetCfg.From, mailjetCfg.Name)

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
