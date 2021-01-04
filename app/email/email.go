package email

import (
	"fmt"

	"go.uber.org/zap"

	"net/smtp"
)

// AuthEmail - struct to authenticate user user
type AuthEmail struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	zapLog   *zap.Logger
}

// Email - interface.
type Email interface {
	SendEmail(to []string, body []byte) error
}

// NewAuthEmail - returns NewAuthEmail object.
func NewAuthEmail(emailFrom string, emailPassword string, smtpHost string, smtpPort string, zapLog *zap.Logger) *AuthEmail {
	return &AuthEmail{emailFrom, emailPassword, smtpHost, smtpPort, zapLog}
}

// SendEmail - send an email using auth.
func (a *AuthEmail) SendEmail(to []string, body []byte) error {
	// Authentication.
	auth := smtp.PlainAuth("", a.from, a.password, a.smtpHost)

	// Sending email.
	if err := smtp.SendMail(a.smtpHost+":"+a.smtpPort, auth, a.from, to, body); err != nil {
		a.zapLog.Error(fmt.Sprintf("Could not send email with err: %v", err))
		return err
	}
	return nil
}
