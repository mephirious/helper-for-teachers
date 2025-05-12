package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET_KEY"))
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "hardwarerump04@gmail.com",
				Name:  "Mailjet Pilot",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "230810@astanait.edu.kz",
					Name:  "passenger 1",
				},
			},
			Subject:  "Your email flight plan!",
			TextPart: "Dear passenger 1, welcome to Mailjet! May the delivery force be with you!",
			HTMLPart: "<h3>Dear passenger 1, welcome to <a href=\"https://www.mailjet.com/\">Mailjet</a>!</h3><br />May the delivery force be with you!",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}
