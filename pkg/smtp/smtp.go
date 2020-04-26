package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
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
	}
}

type Email struct {
	Username, password, Host, From string
	Port                           int
	Timeout                        time.Duration
}

func (e *Email) HasPassword() bool {
	return e.password != ""
}

func (e *Email) Send(email *email.Email) error {
	log.Printf("sending email to %s", strings.Join(email.ReadReceipt, "r"))
	email.From = e.From
	return email.Send(fmt.Sprintf("%s:%d", e.Host, e.Port), smtp.PlainAuth("", e.Username, e.password, e.Host))
}
