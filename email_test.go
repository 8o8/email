package email_test

import (
	"log"
	"testing"

	"github.com/8o8/email"
	"github.com/gobuffalo/envy"
)

var AWS_REGION = ""
var AWS_ACCESS_KEY_ID = ""
var AWS_SECRET_ACCESS_KEY = ""
var MAILGUN_DOMAIN = ""
var MAILGUN_API_KEY = ""
var SENDGRID_API_KEY = ""

// Note for Amazon SES: If account is sandbox, both from and to emails MUST be verified
const (
	senderName     = "8o8/email"
	senderEmail    = "michael.donnici@csanz.edu.au"
	recipientName  = "Mike Donnici"
	recipientEmail = "michael.donnici@gmail.com"
)

func TestAll(t *testing.T) {

	setup()

	t.Run("email", func(t *testing.T) {
		t.Run("testSES", testSES)
		t.Run("testMailgun", testMailgun)
		//t.Run("testSendgrid", testSendgrid)
	})

}

func setup() {
	envy.Load()
	var err error
	AWS_REGION, err = envy.MustGet("AWS_REGION")
	if err != nil {
		log.Fatalln("could not get env var AWS_REGION")
	}
	AWS_ACCESS_KEY_ID, err = envy.MustGet("AWS_ACCESS_KEY_ID")
	if err != nil {
		log.Fatalln("could not get env var AWS_ACCESS_KEY_ID")
	}
	AWS_SECRET_ACCESS_KEY, err = envy.MustGet("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		log.Fatalln("could not get env var AWS_SECRET_ACCESS_KEY")
	}
	MAILGUN_DOMAIN, err = envy.MustGet("MAILGUN_DOMAIN")
	if err != nil {
		log.Fatalln("could not get env var MAILGUN_DOMAIN")
	}
	MAILGUN_API_KEY, err = envy.MustGet("MAILGUN_API_KEY")
	if err != nil {
		log.Fatalln("could not get env var MAILGUN_API_KEY")
	}
	SENDGRID_API_KEY, err = envy.MustGet("SENDGRID_API_KEY")
	if err != nil {
		log.Fatalln("could not get env var SENDGRID_API_KEY")
	}
}

func testMailgun(t *testing.T) {

	cfg := email.MailgunCfg{
		Domain: MAILGUN_DOMAIN,
		APIKey: MAILGUN_API_KEY,
	}
	mg := email.NewMailgun(cfg)

	email := email.Email{
		FromName:     senderName,
		FromEmail:    senderEmail,
		ToName:       recipientName,
		ToEmail:      recipientEmail,
		Subject:      "Mailgun Test",
		PlainContent: "This is the plain text",
		HTMLContent:  "<h1>This is HTML</h1>",
	}

	err := mg.Send(email)
	if err != nil {
		t.Errorf("Send() err = %s", err)
	}
}

func testSendgrid(t *testing.T) {

	cfg := email.SendgridCfg{APIKey: SENDGRID_API_KEY}
	sg := email.NewSendgrid(cfg)

	email := email.Email{
		FromName:     senderName,
		FromEmail:    senderEmail,
		ToName:       recipientName,
		ToEmail:      recipientEmail,
		Subject:      "Sendgrid Test",
		PlainContent: "This is the plain text",
		HTMLContent:  "<h1>This is HTML</h1>",
	}

	err := sg.Send(email)
	if err != nil {
		t.Errorf("Send() err = %s", err)
	}
}

func testSES(t *testing.T) {

	cfg := email.SESCfg{
		AWSRegion: AWS_REGION,
		AWSAccessKeyID: AWS_ACCESS_KEY_ID,
		AWSSecretAccessKey: AWS_SECRET_ACCESS_KEY,
	}
	ses, err := email.NewSES(cfg)
	if err != nil {
		t.Fatalf("email.NewSES() err = %s", err)
	}

	email := email.Email{
		FromName:     senderName,
		FromEmail:    senderEmail,
		ToName:       recipientName,
		ToEmail:      recipientEmail,
		Subject:      "AWS SES Test",
		PlainContent: "This is the plain text",
		HTMLContent:  "<h1>This is HTML</h1>",
	}

	err = ses.Send(email)
	if err != nil {
		t.Errorf("Send() err = %s", err)
	}
}