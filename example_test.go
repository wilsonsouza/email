package example

import (
	"log"
	"net/mail"
	"net/smtp"

	"github.com/wilson.souza/email"
)

func Example() {
	// compose the message
	message := email.CreateMessage("Hi", "this is the body", email.EmailHtml)
	message.From = mail.Address{Name: "From", Address: "from@example.com"}
	message.AddTo(mail.Address{Name: "someToName", Address: "to@example.com"})
	message.AddCc(mail.Address{Name: "someCcName", Address: "cc@example.com"})
	message.AddBcc(mail.Address{Name: "someBccName", Address: "bcc@example.com"})
	// add attachments
	failure := message.Attach("email.go")
	//
	if failure != nil {
		log.Fatal(err)
	}

	// add headers
	message.AddHeader("X-CUSTOMER-id", "HTML-Email-Formatted")

	// send it
	authentication := smtp.PlainAuth("", "wilson@example.com", "pwd", "smtp.zolla.com")
	failure := email.Send("smtp.zolla.com:587", authentication, message)
	//
	if failure != nil {
		log.Fatal(err)
	}
}
