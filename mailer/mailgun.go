package mailer

import (
	"fmt"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// MailGunMailer is the Mailgun specific implementation for the Mailer service
type MailGunMailer struct {
	privateAPIKey string
	publicAPIKey  string
	senderDomain  string
}

// NewMailgunMailer construcs a new MailgunMailer which implements the Mailer interface
func NewMailgunMailer(senderDomain string, privateAPIKey string, publicAPIKey string) MailGunMailer {
	return MailGunMailer{
		senderDomain:  senderDomain,
		privateAPIKey: privateAPIKey,
		publicAPIKey:  publicAPIKey,
	}
}

// SendMail sends out the email via the mailgun SDK
func (mailer MailGunMailer) SendMail(mail Email) (string, error) {
	requestedSenderDomain, err := mail.GetSenderDomain()
	if err != nil {
		return "", ErrInvalidDomain
	}
	if mailer.senderDomain != requestedSenderDomain {
		return "", ErrInvalidDomain
	}
	mg := mailgun.NewMailgun(mailer.senderDomain, mailer.privateAPIKey, mailer.publicAPIKey)
	message := mg.NewMessage(
		mail.From,
		mail.Subject,
		mail.Message,
		mail.To,
	)
	resp, id, err := mg.Send(message)
	if err != nil {
		return "", err
	}
	returnString := fmt.Sprintf("ID:%s, Response:%s", id, resp)
	return returnString, nil
}
