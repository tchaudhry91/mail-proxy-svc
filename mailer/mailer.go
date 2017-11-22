package mailer

import (
	"errors"
	"strings"
)

// Email struct provides a structure with the bare minimum requirements to send an email
type Email struct {
	From    string
	Subject string
	Message string
	To      string
}

var ErrInvalidDomain = errors.New("Invalid domain supplied")

func (mail *Email) GetSenderDomain() (string, error) {
	var senderDomain string
	atIndex := strings.LastIndex(mail.From, "@")
	if atIndex == -1 {
		return senderDomain, ErrInvalidDomain
	}
	senderDomain = mail.From[atIndex+1:]
	return senderDomain, nil
}

type Mailer interface {
	SendMail(mail Email) (err error, response string)
}
