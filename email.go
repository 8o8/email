// Package email is a wrapper for third-party email services such as Sendgrid, Mailgun and Amazon SES
package email

// A sender can send an email
type sender interface {
	Send(Email) error
}

// Email represents an email
type Email struct {
	FromName     string
	FromEmail    string
	ToEmail      string
	ToName       string
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
