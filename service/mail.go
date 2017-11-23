package service

import (
	"context"
	"github.com/tchaudhry91/mail-proxy-svc/mailer"
)

type MailService interface {
	SendEmail(ctx context.Context, From string, Subject string, Message string, To string) (string, error)
}

type mailService struct {
	primaryProvider string
	mailProvider    mailer.Mailer
	config          Configuration
}

type Configuration struct {
	senderDomain  string
	privateApiKey string
	publicApiKey  string
}

func (svc *mailService) InitializePrimaryProvider() {
	switch svc.primaryProvider {
	case "mailgun":
		svc.mailProvider = mailer.NewMailgunMailer(
			svc.config.senderDomain,
			svc.config.privateApiKey,
			svc.config.publicApiKey,
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
