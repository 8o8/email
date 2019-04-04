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
	fmt.Println("Send email via SES...")

	const charset = "UTF-8"

	// Assemble the email.
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
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	_, err := ss.sender.SendEmail(input)

	// Display error messages if they occur.
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
