package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tchaudhry91/mail-proxy-svc/mailer"
	"net/http"
)

var (
	//ErrJSONUnmarshall indicates a bad request where json mashalling fails
	ErrJSONUnmarshall = errors.New("failed to parse json")

	//ErrInvalidInput indicates a bad request where supplied inputs are invalid
	ErrInvalidInput = errors.New("bad input, please recheck")
)

// MakeSendEmailHandler populates the http.Handler with the go-kit endpoint routes
func MakeSendEmailHandler(e endpoint.Endpoint, r *mux.Router) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	sendEmailHandler := httptransport.NewServer(
		e,
		decodeSendEmailRequest,
		encodeSendEmailResponse,
		options...,
	)

	r.Methods("POST").Path("/send").Handler(sendEmailHandler)
}

func decodeSendEmailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request sendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return request, nil
}

func encodeSendEmailResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.WriteHeader(codeFrom(err))
	e := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
	if e != nil {
		panic("Error encoding error, programmer error")
	}
}

func codeFrom(err error) int {
	switch err {
	case ErrJSONUnmarshall:
		return http.StatusBadRequest
	case ErrInvalidInput:
		return http.StatusBadRequest
	case mailer.ErrInvalidDomain:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
