package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/tchaudhry91/mail-proxy-svc/service"
	"net/http"
	"os"
	"strings"
)

// ErrInvalidConfiguration is thrown when incorrect startup parameters have been supplied
var ErrInvalidConfiguration = errors.New("Invalid startup parameters")

// ErrUnsupportedProvider is thrown when a non-standard provider is supplied
var ErrUnsupportedProvider = errors.New("Unsupported Provider")

func main() {
	var (
		primaryProvider string
		senderDomain    string
		privateAPIKey   string
		publicAPIKey    string
		serverAddr      string
	)
	flag.StringVar(
		&serverAddr,
		"serverAddr",
		":9999",
		"Fully qualified server address. Defaults to :9999",
	)
	flag.StringVar(
		&primaryProvider,
		"primaryProvider",
		"",
		"Name of the email provider: One of [mailgun,]",
	)
	flag.StringVar(
		&senderDomain,
		"senderDomain",
		"",
		"Sender domain registered with the provider",
	)
	flag.StringVar(
		&privateAPIKey,
		"privateAPIKey",
		"",
		"PrivateAPI key for the selected provider",
	)
	flag.StringVar(
		&publicAPIKey,
		"publicAPIKey",
		"",
		"PublicAPIKey for the selected provider if required",
	)
	flag.Parse()
	primaryProvider = strings.ToLower(primaryProvider)
	config := service.Configuration{
		SenderDomain:  senderDomain,
		PrivateAPIKey: privateAPIKey,
		PublicAPIKey:  privateAPIKey,
	}
	err := validateConfig(primaryProvider, config)
	if err != nil {
		fmt.Printf("Encountered Err: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}
	logger := log.NewJSONLogger(os.Stderr)
	router := mux.NewRouter()

	svc := service.NewMailService(primaryProvider, config)
	svc = service.NewLoggingMiddleware(logger, svc)
	endpoint := service.MakeSendEmailEndpoint(svc)
	service.MakeSendEmailHandler(endpoint, router)

	http.Handle("/", router)

	logger.Log("msg", "HTTP", "addr", serverAddr)
	logger.Log("err", http.ListenAndServe(serverAddr, nil))
}

func validateConfig(primaryProvider string, config service.Configuration) error {
	switch primaryProvider {
	case "mailgun":
		if config.SenderDomain != "" && config.PrivateAPIKey != "" && config.PublicAPIKey != "" {
			return nil
		}
		return ErrInvalidConfiguration
	default:
		return ErrUnsupportedProvider
	}
}
