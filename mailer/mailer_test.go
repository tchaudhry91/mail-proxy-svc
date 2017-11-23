package mailer

import (
	"testing"
)

var samples []Email
var correspondingSampleDomains []string

func init() {
	samples = []Email{
		{
			From:    "tanmay.chaudhry@gmail.com",
			Subject: "Test Email",
			Message: "No message",
			To:      "xyz@gmail.com",
		},
		{
			From:    "tanmay@tux-sudo.com",
			Subject: "Test Email",
			Message: "No message",
			To:      "xyz@gmail.com",
		},
		{
			From:    "1231111333",
			Subject: "Fooling you",
			Message: "No message",
			To:      "send@sss.com",
		},
	}
	correspondingSampleDomains = []string{
		"gmail.com",
		"tux-sudo.com",
		"",
	}
}

func TestSenderDomain(t *testing.T) {
	for index, mail := range samples {
		want := correspondingSampleDomains[index]
		have, err := mail.GetSenderDomain()
		if err != nil {
			if want != "" {
				t.Errorf("%s : want %q, have %q", mail.From, want, have)
			}
		}
		if want != have {
			t.Errorf("%s : want %q, have %q", mail.From, want, have)
		}
	}
}
