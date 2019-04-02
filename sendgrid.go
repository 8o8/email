package email

import (
	"fmt"
	"net/http"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Sendgrid implements the sender interface
type Sendgrid struct {
	cfg    SendgridCfg
	sender *sendgrid.Client
}

// SendgridCfg defines the configuration for Sendgrid
type SendgridCfg struct {
	APIKey string
}

// NewSendgrid returns a pointer to a SendgridEmail
func NewSendgrid(cfg SendgridCfg) *Sendgrid {
	sg := sendgrid.NewSendClient(cfg.APIKey)
	return &Sendgrid{
		cfg:    cfg,
		sender: sg,
	}
}

// Send sends an email
func (sg *Sendgrid) Send(e Email) error {
	fmt.Println("Send email via sendgrid...")

	message := prepare(e)

	for _, a := range e.Attachments {
		attach(a, message)
	}

	response, err := sg.sender.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid.Client.Send() err = %s", err)
	}

	// error if non-200 response code
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("sendgrid.Client.Send() response (%d) - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	return nil
}

func prepare(e Email) *mail.SGMailV3 {
	from := mail.NewEmail(e.FromName, e.FromEmail)
	subject := e.Subject
	to := mail.NewEmail(e.ToName, e.ToEmail)
	plainTextContent := e.PlainContent
	htmlContent := e.HTMLContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	return message
}

func attach(a Attachment, message *mail.SGMailV3) {
	if a.isOK() {
		message.AddAttachment(newAttachment(a))
	}
}

func (a Attachment) isOK() bool {
	return a.Base64Content != "" && a.FileName != "" && a.MIMEType != ""
}

func newAttachment(a Attachment) *mail.Attachment {
	ma := mail.NewAttachment()
	ma.SetContent(a.Base64Content)
	ma.SetType(a.MIMEType)
	ma.SetFilename(a.FileName)
	ma.SetDisposition("attachment") // no "inline" for now
	// ma.SetContentID("Attachment...") // used for inline attachments
	return ma
}
