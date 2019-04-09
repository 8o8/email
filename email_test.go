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
var TEST_SENDER_NAME = ""
var TEST_SENDER_EMAIL = ""
var TEST_RECIPIENT_NAME = ""
var TEST_RECIPIENT_EMAIL = ""
var TEST_SUBJECT = ""
var TEST_HTML_CONTENT = ""
var TEST_PLAIN_CONTENT = ""

// By default, no tests will run
var TEST_MAILGUN = false
var TEST_SENDGRID = false
var TEST_SES = false

func TestAll(t *testing.T) {

	setup()

	t.Run("email", func(t *testing.T) {
		t.Run("testSES", testSES)
		t.Run("testMailgun", testMailgun)
		t.Run("testSendgrid", testSendgrid)
	})
}

func setup() {

	envy.Load()

	// Optional, but one should probably be configured for anything to happen!
	AWS_REGION = envy.Get("AWS_REGION", "")
	AWS_ACCESS_KEY_ID = envy.Get("AWS_ACCESS_KEY_ID", "")
	AWS_SECRET_ACCESS_KEY = envy.Get("AWS_SECRET_ACCESS_KEY", "")
	MAILGUN_DOMAIN = envy.Get("MAILGUN_DOMAIN", "")
	MAILGUN_API_KEY = envy.Get("MAILGUN_API_KEY", "")
	SENDGRID_API_KEY = envy.Get("SENDGRID_API_KEY", "")

	// Required
	var err error
	TEST_SENDER_NAME, err = envy.MustGet("TEST_SENDER_NAME")
	if err != nil {
		log.Fatalf("setup() could not get env var TEST_SENDER_NAME")
	}
	TEST_SENDER_EMAIL, err = envy.MustGet("TEST_SENDER_EMAIL")
	if err != nil {
		log.Fatalf("setup() could not get env var TEST_SENDER_EMAIL")
	}
	TEST_RECIPIENT_NAME, err = envy.MustGet("TEST_RECIPIENT_NAME")
	if err != nil {
		log.Fatalf("setup() could not get env var TEST_RECIPIENT_NAME")
	}
	TEST_RECIPIENT_EMAIL, err = envy.MustGet("TEST_RECIPIENT_EMAIL")
	if err != nil {
		log.Fatalf("setup() could not get env var TEST_RECIPIENT_EMAIL")
	}

	// test flags
	test_mailgun := envy.Get("TEST_MAILGUN", "false")
	if test_mailgun == "true" {
		TEST_MAILGUN = true
	}
	test_sendgrid := envy.Get("TEST_SENDGRID", "false")
	if test_sendgrid == "true" {
		TEST_SENDGRID = true
	}
	test_ses := envy.Get("TEST_SES", "false")
	if test_ses == "true" {
		TEST_SES = true
	}
}

func testMailgun(t *testing.T) {

	if !TEST_MAILGUN {
		t.Log("TEST_MAILGUN = false")
		return
	}

	cfg := email.MailgunCfg{
		Domain: MAILGUN_DOMAIN,
		APIKey: MAILGUN_API_KEY,
	}
	mg := email.NewMailgun(cfg)

	email := email.Email{
		FromName:     TEST_SENDER_NAME,
		FromEmail:    TEST_SENDER_EMAIL,
		ToName:       TEST_RECIPIENT_NAME,
		ToEmail:      TEST_RECIPIENT_EMAIL,
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

	if !TEST_SENDGRID {
		t.Log("TEST_SENDGRID = false")
		return
	}

	cfg := email.SendgridCfg{APIKey: SENDGRID_API_KEY}
	sg := email.NewSendgrid(cfg)

	email := email.Email{
		FromName:     TEST_SENDER_NAME,
		FromEmail:    TEST_SENDER_EMAIL,
		ToName:       TEST_RECIPIENT_NAME,
		ToEmail:      TEST_RECIPIENT_EMAIL,
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

	if !TEST_SES {
		t.Log("TEST_SES = false")
		return
	}

	cfg := email.SESCfg{
		AWSRegion:          AWS_REGION,
		AWSAccessKeyID:     AWS_ACCESS_KEY_ID,
		AWSSecretAccessKey: AWS_SECRET_ACCESS_KEY,
	}
	ses, err := email.NewSES(cfg)
	if err != nil {
		t.Fatalf("email.NewSES() err = %s", err)
	}

	email := email.Email{
		FromName:     TEST_SENDER_NAME,
		FromEmail:    TEST_SENDER_EMAIL,
		ToName:       TEST_RECIPIENT_NAME,
		ToEmail:      TEST_RECIPIENT_EMAIL,
		Subject:      "AWS SES Test",
		PlainContent: "This is the plain text",
		HTMLContent:  "<h1>This is HTML</h1>",
	}

	err = ses.Send(email)
	if err != nil {
		t.Errorf("Send() err = %s", err)
	}
}
