package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

// Mailgun implements the sender interface
type Mailgun struct {
	cfg    MailgunCfg
	sender mailgun.Mailgun
}

// MailgunCfg defines the configuration for Mailgun
type MailgunCfg struct {
	Domain string
	APIKey string
}

// NewMailgun returns a pointer to a local Mailgun value
func NewMailgun(cfg MailgunCfg) *Mailgun {
	mg := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)
	return &Mailgun{
		cfg:    cfg,
		sender: mg,
	}
}

// Send sends an email
func (mg *Mailgun) Send(e Email) error {
	fmt.Println("Send email via mailgun...")

	// The message object allows you to add attachments and Bcc recipients
	message := mg.sender.NewMessage(e.From(), e.Subject, e.PlainContent, e.To())
	message.SetHtml(e.HTMLContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.sender.Send(ctx, message)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	// error if non-200 response code
	// if response.StatusCode < 200 || response.StatusCode > 299 {
	// 	return fmt.Errorf("sendgrid.Client.Send() response (%d) - %s", response.StatusCode, http.StatusText(response.StatusCode))
	// }

	return nil
}
