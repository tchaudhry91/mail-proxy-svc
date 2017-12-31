package service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"strconv"
	"time"
)

// InstrumentingMiddleware is a middleware to produce metrics for the mail-proxy-service
type InstrumentingMiddleware struct {
	emailCount      metrics.Counter
	requestDuration metrics.Histogram
	next            MailService
}

// SendEmail is a method wrapper around the internal SendEmail to instrument the request
func (mw InstrumentingMiddleware) SendEmail(ctx context.Context, from string, subject string, message string, to string, html bool) (response string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SendEmail",
			"error", fmt.Sprint(err != nil),
			"to", to,
			"from", from,
			"subject", subject,
			"html", strconv.FormatBool(html),
		}
		mw.emailCount.With(lvs...).Add(1)
		mw.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	response, err = mw.next.SendEmail(ctx, from, subject, message, to, html)
	return
}

// NewInstrumentingMiddleware is a constructor for the instrumenting middleware used by the Mail-Proxy-Service
func NewInstrumentingMiddleware(emailCount metrics.Counter, requestDuration metrics.Histogram, next MailService) InstrumentingMiddleware {
	return InstrumentingMiddleware{emailCount, requestDuration, next}
}
