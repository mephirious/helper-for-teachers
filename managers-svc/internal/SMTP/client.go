package smtp

import (
	"fmt"
	"net/mail"

	gomail "gopkg.in/mail.v2"
)

type Attachment struct {
	Filename string
	Data     []byte
}

type Client struct {
	host     string
	port     int
	email    string
	password string
}

func NewClient(host string, port int, email, password string) *Client {
	return &Client{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (c *Client) Send(to []string, subject, body string) error {
	for _, email := range to {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf("invalid email address %q: %w", email, err)
		}
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", c.email)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	dialer := gomail.NewDialer(c.host, c.port, c.email, c.password)

	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("sending failed: %w", err)
	}

	return nil
}
