package smtp

import (
	"log"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

type Email struct {
	Username, password, Host, From string
	Port                           int
	Timeout                        time.Duration
	sendFunc                       func(string, *email.Email) error
}

func (e *Email) HasPassword() bool {
	return e.password != ""
}

func (e *Email) Send(email *email.Email) error {
	log.Printf("sending email to %s", strings.Join(email.ReadReceipt, "r"))
	return e.sendFunc(e.From, email)
}
