package service

import (
	"context"
	"github.com/tchaudhry91/mail-proxy-svc/mailer"
)

// MailService is a generic mail proxying service which can send email via various backends
type MailService interface {
	SendEmail(ctx context.Context, from string, subject string, message string, to string) (string, error)
}

type mailService struct {
	primaryProvider string
	mailProvider    mailer.Mailer
	config          Configuration
}

// Configuration is a struct to maintain provider config in a single place
type Configuration struct {
	SenderDomain  string
	PrivateAPIKey string
	PublicAPIKey  string
}

func (svc *mailService) InitializePrimaryProvider() {
	switch svc.primaryProvider {
	case "mailgun":
		svc.mailProvider = mailer.NewMailgunMailer(
			svc.config.SenderDomain,
			svc.config.PrivateAPIKey,
			svc.config.PublicAPIKey,
		)
	}
}

// NewMailService is a constructor for the mailService
func NewMailService(primaryProvider string, config Configuration) MailService {
	svc := mailService{
		primaryProvider: primaryProvider,
		config:          config,
	}
	svc.InitializePrimaryProvider()
	return svc
}

func (svc mailService) SendEmail(ctx context.Context, from string, subject string, message string, to string) (string, error) {
	email := mailer.Email{
		From:    from,
		Subject: subject,
		Message: message,
		To:      to,
	}
	resp, err := svc.mailProvider.SendMail(email)
	return resp, err
}
