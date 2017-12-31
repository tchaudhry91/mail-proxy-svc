package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

// LoggingMiddleware is a middleware to produce structured logs for the mail-proxy-service
type LoggingMiddleware struct {
	logger log.Logger
	next   MailService
}

// SendEmail is a method wrapper around the internal SendEmail to log the request
func (mw LoggingMiddleware) SendEmail(ctx context.Context, from string, subject string, message string, to string, html bool) (response string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SendEmail",
			"from", from,
			"subject", subject,
			"to", to,
			"html", html,
			"response", response,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	response, err = mw.next.SendEmail(ctx, from, subject, message, to, html)
	return
}

// NewLoggingMiddleware is a constructor for the logging middleware used by the MailService
func NewLoggingMiddleware(logger log.Logger, next MailService) LoggingMiddleware {
	return LoggingMiddleware{logger, next}
}
