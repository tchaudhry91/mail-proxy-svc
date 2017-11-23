package mailer

import (
	"fmt"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type MailGunMailer struct {
	privateApiKey string
	publicApiKey  string
	senderDomain  string
}

func NewMailgunMailer(senderDomain string, privateApiKey string, publicApiKey string) MailGunMailer {
	return MailGunMailer{
		senderDomain:  senderDomain,
		privateApiKey: privateApiKey,
		publicApiKey:  publicApiKey,
	}
}

func (mailer MailGunMailer) SendMail(mail Email) (string, error) {
	requestedSenderDomain, err := mail.GetSenderDomain()
	if err != nil {
		return "", ErrInvalidDomain
	}
	if mailer.senderDomain != requestedSenderDomain {
		return "", ErrInvalidDomain
	}
	mg := mailgun.NewMailgun(mailer.senderDomain, mailer.privateApiKey, mailer.publicApiKey)
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
