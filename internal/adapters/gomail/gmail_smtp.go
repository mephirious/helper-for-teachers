package gomail

import (
	"github.com/mephirious/helper-for-teachers/services/auth-svc/internal/domain"
	gomailpkg "github.com/mephirious/helper-for-teachers/services/auth-svc/pkg/gomail"
)

func NewGomailService(sender *gomailpkg.Sender) domain.EmailSender {
	return sender
}
