// Package email is a wrapper for third-party email services - Sendgrid, Mailgun and Amazon SES
package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	//"github.com/aws/aws-sdk-go/service/ses"
	gomail "gopkg.in/gomail.v2"
)

// A sender can send email
type sender interface {
	Send(Email) error
}

// Email represents an email
type Email struct {
	FromName     string
	FromEmail    string
	ToName       string
	ToEmail      string
	Subject      string
	PlainContent string
	HTMLContent  string
	Attachments  []Attachment
}

// Attachment to an email
type Attachment struct {
	MIMEType      string
	FileName      string
	Base64Content string
}

// Send the email using the specified sender
func (e Email) Send(srvc sender) error {
	return srvc.Send(e)
}

// To returns an appropriately formatted recipient string
func (e Email) To() string {
	return nameEmail(e.ToName, e.ToEmail)
}

// From returns an appropriately formatted sender string
func (e Email) From() string {
	return nameEmail(e.FromName, e.FromEmail)
}

// Raw converts the Email value to a raw message string in line with RFC 5322
// This is required for sending emails with attachments via Amazon SES.
func (e Email) Raw() (string, error) {

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", e.FromName, e.FromEmail))
	msg.SetHeader("To", fmt.Sprintf("%s <%s>", e.ToName, e.ToEmail))
	msg.SetHeader("Subject", e.Subject)
	msg.SetBody("text/plain", e.PlainContent)
	msg.AddAlternative("text/html", e.HTMLContent)

	for i := 0; i < len(e.Attachments); i++ {

		xb, err := base64.StdEncoding.DecodeString(e.Attachments[i].Base64Content)
		if err != nil {
			return "", fmt.Errorf("DecodeString() err = %s", err)
		}
		msg.Attach(e.Attachments[i].FileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(xb)
			return err
		}))
	}

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	return string(emailRaw.Bytes()), nil
}

// nameEmail returns a string in the format: 'Name <name@somewhere.com>' if name
// is specified. Otherwise, it simply returns the email address.
func nameEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
