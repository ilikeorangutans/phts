package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

func NewEmailSender(host string, port int, username, password string) *Email {
	return &Email{
		username: username,
		password: password,
		host:     host,
		port:     port,
		timeout:  60 * time.Second,
	}
}

type Email struct {
	username, password, host string
	port                     int
	timeout                  time.Duration
}

func (e *Email) Send(email *email.Email) error {
	log.Printf("sending email to %s", strings.Join(email.ReadReceipt, "r"))
	return email.Send(fmt.Sprintf("%s:%d", e.host, e.port), smtp.PlainAuth("", e.username, e.password, e.host))
}
