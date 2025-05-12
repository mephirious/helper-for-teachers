package mailjet

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type MailjetClient struct {
	client *mailjet.Client
	from   string
	name   string
}

func NewMailjetClient(client *mailjet.Client, from string, name string) *MailjetClient {
	return &MailjetClient{
		client: client,
		from:   from,
		name:   name,
	}
}

func (m *MailjetClient) SendEmail(toEmail, toName, subject, textBody, htmlBody string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.from,
				Name:  m.name,
			},
			To: &mailjet.RecipientsV31{
				{
					Email: toEmail,
					Name:  toName,
				},
			},
			Subject:  subject,
			TextPart: textBody,
			HTMLPart: htmlBody,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
