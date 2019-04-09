# email for Go

A basic wrapper for Amazon SES, Mailgun and Sendgrid APIs

## Installation

```
go get -u github.com/8o8/email
```

## Tests Require Env Vars

To run the tests set up a `.env` file with the following values:

> **Don't forget to add `.env` to `.gitignore`**. If you accidentally push the `.env` file to a public repo, the keys will be detected, alarms will ring and Mailgun/Sendgrid will immediately suspend your account.

```env
# For Amazon SES - also need IAM, Policy and all that guff
AWS_REGION="amazon-ses-region"
AWS_ACCESS_KEY_ID="your-amazon-access-key-id"
AWS_SECRET_ACCESS_KEY="your-amazon-secret-access-key"

# Mailgun is easier
MAILGUN_DOMAIN="your.mailgun.domain"
MAILGUN_API_KEY="your-mailgun-api-key"

# Sengrid, easier again
SENDGRID_API_KEY="your-sendgrid-api-key"

# If Amazon SES account is still in sandbox mode, these email addresses will
# need to be verified before they will work
TEST_SENDER_NAME="8o8 Test Mail"
TEST_SENDER_EMAIL="verified@email.oi"
TEST_RECIPIENT_NAME="Your Name"
TEST_RECIPIENT_EMAIL="youremail@moomoo.com"

# Skip test for any value other than "true"
TEST_MAILGUN="true"
TEST_SENDGRID="" # false
TEST_SES=true="" # false
```

## Usage

```go
package main

import "github.com/8o8/email"

func main() {

    em := email.Email{
        FromName: "Send Name",
        FromEmail: "sender@email.oi",
        ToName: "Recipient Name",
        ToEmail: "recipient@email.oi",
        Subject: "The Subject",
        PlainContent: "This is the plain text",
        HTMLContent: "<h1>This is HTML</h1>",
    }

    // cfg for the specific api
    cfg := email.MailgunCfg{
        APIKey: apiKey,
        Domain: domain,
    }
    mx := email.NewMailgun(cfg)

    err := mx.Send(em)
    if err != nil {
        fmt.Println(err)
    }
}
```

## Version

Current: v0.1.0

As stated in item 4 of the [SemVer specification](<https://semver.org/>):

>Major version zero (0.y.z) is for initial development. Anything may change at
any time. The public API should not be considered stable.

