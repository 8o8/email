// Package email is a wrapper for third-party email services - Sendgrid, Mailgun and Amazon SES
package email

import "fmt"

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

// nameEmail returns a string in the format: 'Name <name@somewhere.com>' if name
// is specified. Otherwise, it simply returns the email address.
func nameEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
