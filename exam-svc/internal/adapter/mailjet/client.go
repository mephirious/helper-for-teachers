package mailjet

import (
	"fmt"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type MailjetClient struct {
	client *mailjet.Client
	from   string
	name   string
}

func NewMailjetClient(cfg *config.MailjetConfig) *MailjetClient {
	client := mailjet.NewMailjetClient(cfg.API, cfg.KEY)

	return &MailjetClient{
		client: client,
		from:   cfg.From,
		name:   cfg.Name,
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
