package service

import (
	"context"
	"github.com/tchaudhry91/mail-proxy-svc/mailer"
)

// MailService is a generic mail proxying service which can send email via various backends
type MailService interface {
	SendEmail(ctx context.Context, From string, Subject string, Message string, To string) (string, error)
}

type mailService struct {
	primaryProvider string
	mailProvider    mailer.Mailer
	config          Configuration
}

// Configuration is a struct to maintain provider config in a single place
type Configuration struct {
	senderDomain  string
	privateAPIKey string
	publicAPIKey  string
}

func (svc *mailService) InitializePrimaryProvider() {
	switch svc.primaryProvider {
	case "mailgun":
		svc.mailProvider = mailer.NewMailgunMailer(
			svc.config.senderDomain,
			svc.config.privateAPIKey,
			svc.config.publicAPIKey,
		)
	}
}

func (svc *mailService) SendEmail(ctx context.Context, From string, Subject string, Message string, To string) (string, error) {
	email := mailer.Email{
		From:    From,
		Subject: Subject,
		Message: Message,
		To:      To,
	}
	resp, err := svc.mailProvider.SendMail(email)
	return resp, err
}
