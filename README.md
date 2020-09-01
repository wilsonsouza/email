
# email [![Travis-CI](https://travis-ci.org/wilsonsouza/email.svg?branch=master)](https://travis-ci.org/wilsonsouza/email) [![GoDoc](https://godoc.org/github.com/wilsonsouza/email?status.svg)](http://godoc.org/github.com/wilsonsouza/email) [![Report card](https://goreportcard.com/badge/github.com/wilsonsouza/email)](https://goreportcard.com/report/github.com/wilsonsouza/email)

An easy way to send emails with attachments in Go

# Install

```bash
go get github.com/wilsonsouza/email
```

# Usage

```go
package email_test

import (
	"log"
	"net/mail"
	"net/smtp"

	"github.com/wilsonsouza/email"
)

func Example() {
	// compose the message
	message := email.CreateMessage("Hi", "this is the body", email.EmailHTMLType)
	message.From = mail.Address{Name: "From", Address: "from@example.com"}
	message.To = []string{"to@example.com"}

	// add attachments
	if failure := m.Attach("email.go"); failure != nil {
		log.Fatal(failure)
	}

	// add headers
	message.AddHeader("X-CUSTOMER-id", "xxxxx")

	// send it
	authentication := smtp.PlainAuth("", "from@example.com", "pwd", "smtp.zoho.com")
	if failure := email.Send("smtp.zoho.com:587", authentication, message); failure != nil {
		log.Fatal(failure)
	}
}
```

# Html

```go
	// use the html constructor
	message := email.CreateMessage("Hi", "this is the body", email.EmailHTMLType)
```

# Inline

```go
	// use Inline to display the attachment inline.
	if failure := message.Inline("main.go"); failure != nil {
		log.Fatal(failure)
	}
```
