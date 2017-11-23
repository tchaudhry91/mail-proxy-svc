package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type sendEmailRequest struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	To      string `json:"to"`
}

type sendEmailResponse struct {
	Response string `json:"response"`
	Error    string `json:"err,omitempty"`
}

// MakeSendEmailEndpoint create a JSON backed go-kit endpoint for mail proxy service
func MakeSendEmailEndpoint(svc MailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sendEmailRequest)
		respString, err := svc.SendEmail(ctx, req.From, req.Subject, req.Message, req.To)
		if err != nil {
			return sendEmailResponse{
				Response: respString,
				Error:    err.Error(),
			}, nil
		}
		return sendEmailResponse{
			Response: respString,
			Error:    "",
		}, nil
	}
}
