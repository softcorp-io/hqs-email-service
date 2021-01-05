package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	email "github.com/softcorp-io/hqs-email-service/email"
	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_email_service"
	"go.uber.org/zap"
)

// Handler - struct used through program and passed to go-micro.
type Handler struct {
	zapLog *zap.Logger
	email  email.Email
}

// NewHandler returns a Handler object
func NewHandler(zapLog *zap.Logger, email email.Email) *Handler {
	return &Handler{zapLog, email}
}

// Ping - used for other service to check connection
func (s *Handler) Ping(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	return &proto.Response{}, nil
}

// SendResetPasswordEmail - sends an email to user with a token that can be used to reset the password.
func (s *Handler) SendResetPasswordEmail(ctx context.Context, email *proto.ResetPasswordEmail) (*proto.Response, error) {
	s.zapLog.Info("Recieved new request")
	// Generate email from HTML
	t, err := template.ParseFiles("template/reset_password.html")
	if err != nil {
		s.zapLog.Error(fmt.Sprintf("Could not get template from file pathwith err: %v", err))
		return nil, err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Hqs Reset Password Request \n%s\n\n", mimeHeaders)))

	// get base link
	resetBaseURL, ok := os.LookupEnv("RESET_PASSWORD_URL")
	if !ok {
		s.zapLog.Error("Could not get RESET_PASSWORD_URL")
		return nil, errors.New("Could not get RESET_PASSWORD_URL")
	}
	link := resetBaseURL + email.Token

	t.Execute(&body, struct {
		Name string
		Link string
	}{
		Name: email.Name,
		Link: link,
	})

	// Sending email.
	err = s.email.SendEmail(email.To, body.Bytes())
	if err != nil {
		s.zapLog.Error(fmt.Sprintf("Could not send email with err: %v", err))
		return nil, err
	}

	s.zapLog.Info("Reset email successfully sent")

	res := &proto.Response{}

	return res, nil
}
