package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	stdprom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tchaudhry91/mail-proxy-svc/service"
	"net/http"
	"os"
	"strings"
	"time"
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

	// Base service
	svc := service.NewMailService(primaryProvider, config)

	// Logging Middleware
	logger := log.NewJSONLogger(os.Stderr)
	svc = service.NewLoggingMiddleware(logger, svc)

	// Prometheus Middleware
	fieldKeys := []string{"method", "error", "to", "from", "subject", "html"}
	emailCount := kitprom.NewCounterFrom(
		stdprom.CounterOpts{
			Namespace: "mail_service",
			Name:      "email_count_total",
			Help:      "Total Emails Requested",
		},
		fieldKeys,
	)
	requestDuration := kitprom.NewSummaryFrom(
		stdprom.SummaryOpts{
			Namespace: "mail_service",
			Name:      "email_processing_seconds",
			Help:      "Time taken per request",
		},
		fieldKeys,
	)
	svc = service.NewInstrumentingMiddleware(emailCount, requestDuration, svc)

	// Go-Kit Endpoint
	endpoint := service.MakeSendEmailEndpoint(svc)

	// Request Router + Transport
	router := mux.NewRouter()
	service.MakeSendEmailHandler(endpoint, router)

	router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	logger.Log("timestamp", time.Now(), "msg", "HTTP Server Working", "addr", serverAddr)
	logger.Log("err", http.ListenAndServe(serverAddr, router))
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
