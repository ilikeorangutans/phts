// +build debug

package smtp

import (
	"log"
	"strings"

	"github.com/jordan-wright/email"
)

func NewEmailSender(host string, port int, username, password, from string) *Email {
	return &Email{
		sendFunc: func(from string, email *email.Email) error {
			log.Printf("Receipient: %s", strings.Join(email.To, ", "))
			log.Printf("Subject: %s", email.Subject)
			log.Printf("%s", email.Text)

			return nil
		},
	}
}
