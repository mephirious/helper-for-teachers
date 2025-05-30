package ports

type SMTP interface {
	Send(to []string, subject, body string) error
}
