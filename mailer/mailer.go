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

// ErrInvalidDomain is returned when an invalid domain is specified for sending
var ErrInvalidDomain = errors.New("Invalid domain supplied")

// GetSenderDomain extracts the sender domain from the From address
func (mail *Email) GetSenderDomain() (string, error) {
	var senderDomain string
	atIndex := strings.LastIndex(mail.From, "@")
	if atIndex == -1 {
		return senderDomain, ErrInvalidDomain
	}
	senderDomain = mail.From[atIndex+1:]
	return senderDomain, nil
}

// Mailer is a generic interface which can be implemented by specific providers
type Mailer interface {
	SendMail(mail Email) (string, error)
}
