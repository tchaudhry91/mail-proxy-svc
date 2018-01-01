# Mail Proxy Service
[![Build Status](https://travis-ci.org/tchaudhry91/mail-proxy-svc.svg?branch=master)](https://travis-ci.org/tchaudhry91/mail-proxy-svc)
[![Go Report Card](https://goreportcard.com/badge/github.com/tchaudhry91/mail-proxy-svc)](https://goreportcard.com/report/github.com/tchaudhry91/mail-proxy-svc)

A simple mail proxying service that proxies mail through popular email providers.

Supported Providers:
- Mailgun

(Hopefully, more to come)

Invoke the binary using the following parameters:

```
Usage of ./mail-proxy-service:
  -primaryProvider string
    	Name of the email provider: One of [mailgun,]
  -privateAPIKey string
    	PrivateAPI key for the selected provider
  -publicAPIKey string
    	PublicAPIKey for the selected provider if required
  -senderDomain string
    	Sender domain registered with the provider
  -serverAddr string
    	Fully qualified server address. Defaults to :9999 (default ":9999")
```

The server will start up with a message like:

`{"addr":":9999","msg":"HTTP Server Working","time":"2018-01-01T12:12:54.143877309+05:30"}`

Send in requests to the "/send" endpoint with the following json structure
```json
{
  "from":"tanmay@mail.tux-sudo.com",
  "to":"sample@tux-sudo.com",
  "subject":"Do you work? Text",
  "message":"Just Testing a Normal Text Email",
  "html":false
}
```

The `from` address should end with a domain with which the service was started : `-senderDomain`

If the `message` contains html, send in `html: true` instead.

Sample request:
`curl -XPOST http://localhost:9999/send -d "@sample_text.json"`

Sample response:
`{"response":"ID:\u003c20180101064922.1.E0B48A1CFFD03A9A@mail.tux-sudo.com\u003e, Response:Queued. Thank you."}`

The rest is up to your provider :)


This is primarily a learning project for go-kit [https://gokit.io/]

Tanmay Chaudhry

tanmay.chaudhry@gmail.com
