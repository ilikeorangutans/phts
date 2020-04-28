// +build !debug

package smtp

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

func NewEmailSender(host string, port int, username, password, from string) *Email {
	return &Email{
		Username: username,
		password: password,
		From:     from,
		Host:     host,
		Port:     port,
		Timeout:  60 * time.Second,
		sendFunc: func(from string, email *email.Email) error {

			email.From = from
			return email.Send(fmt.Sprintf("%s:%d", host, port), smtp.PlainAuth("", username, password, host))
		},
	}
}
