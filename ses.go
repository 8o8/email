package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SES implements the sender interface
type SES struct {
	cfg    SESCfg
	sender *ses.SES
}

// SESCfg provides config info
type SESCfg struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
}

// NewSES returns a pointer to an SES sender
func NewSES(cfg SESCfg) (*SES, error) {
	sndr := &SES{
		cfg: cfg,
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, ""),
	})
	if err != nil {
		return sndr, err
	}
	sndr.sender = ses.New(sess)
	return sndr, nil
}

// Send sends an email
func (ss *SES) Send(e Email) error {

	// If there are attachments use sendRaw().
	if len(e.Attachments) > 0 {
		return ss.sendRaw(e)
	}

	const charset = "UTF-8"
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(e.To()),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(e.HTMLContent),
				},
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(e.PlainContent),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.From()),
	}

	_, err := ss.sender.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return fmt.Errorf("%s - %s", ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return fmt.Errorf("%s - %s", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return fmt.Errorf("%s - %s", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				return fmt.Errorf("%s", aerr.Error())
			}
		}
		return err
	}

	return nil
}

// sendRaw sends a raw (RFC 5322) email - required for attachments via SES
// ref: https://docs.aws.amazon.com/ses/latest/DeveloperGuide/send-email-raw.html
func (ss *SES) sendRaw(e Email) error {

	rawEmail, err := e.Raw()
	if err != nil {
		return fmt.Errorf("Raw() err = %s", err)
	}

	source := aws.String(e.FromEmail)
	destinations := []*string{aws.String(e.ToEmail)}
	message := ses.RawMessage{Data: []byte(rawEmail)}
	input := &ses.SendRawEmailInput{
		Source:       source,
		Destinations: destinations,
		RawMessage:   &message,
	}

	_, err = ss.sender.SendRawEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return fmt.Errorf("%s - %s", ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return fmt.Errorf("%s - %s", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return fmt.Errorf("%s - %s", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				return fmt.Errorf("%s - %s", ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				return fmt.Errorf("%s - %s", ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				return fmt.Errorf("%s", aerr.Error())
			}
		}
		return err
	}

	return nil
}
