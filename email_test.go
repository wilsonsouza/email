package email

import (
	"net/mail"
	"strings"
	"testing"
)

func TestAttachment(test *testing.T) {
	message := CreateMessage("Hi", "this is the body", EmailTextType)
	failure := message.AttachBuffer("test.ics", []byte("test"), false)
	//
	if failure != nil {
		test.Fatal(failure)
	}
	//
	if strings.Contains(string(message.ToBytes()), "text/calendar") == false {
		test.Fatal("Issue with mailer")
	}
}

func TestHeaders(test *testing.T) {
	message := CreateMessage("Hi", "this is the body", EmailTextType)
	message.AddHeader("X-HEADER-KEY", "HEADERVAL")
	//
	if strings.Contains(string(message.ToBytes()), "X-HEADER-KEY: HEADERVAL\r\n") == false {
		test.Fatal("Could not find header in message")
	}
}

func TestAddTo(test *testing.T) {
	message := CreateMessage("Hi", "this is the body", EmailTextType)
	names := []string{"FirstName", "SecondName"}
	addresses := []string{"FirstAddress", "SecondAddress"}
	firstaddress := mail.Address{Name: names[0], Address: addresses[0]}
	message.AddTo(firstaddress)
	//
	if message.To[0] != firstaddress.String() {
		test.Fatal("Incorrect first element")
	}
	//
	secondaddress := mail.Address{Name: names[1], Address: addresses[1]}
	message.AddTo(secondaddress)
	//
	if message.To[1] != secondaddress.String() {
		test.Fatal("Incorrect second element")
	}
}

func TestAddCc(test *testing.T) {
	message := CreateMessage("Hi", "this is the body", EmailTextType)
	names := []string{"FirstName", "SecondName"}
	addresses := []string{"FirstAddress", "SecondAddress"}
	firstaddress := mail.Address{Name: names[0], Address: addresses[0]}
	message.AddCc(firstaddress)
	//
	if message.Cc[0] != firstaddress.String() {
		test.Fatal("Incorrect first element")
	}
	//
	secondaddress := mail.Address{Name: names[1], Address: addresses[1]}
	message.AddCc(secondaddress)
	//
	if message.Cc[1] != secondaddress.String() {
		test.Fatal("Incorrect second element")
	}
}

func TestAddBcc(test *testing.T) {
	message := CreateMessage("Hi", "this is the body", EmailTextType)
	names := []string{"FirstName", "SecondName"}
	addresses := []string{"FirstAddress", "SecondAddress"}
	firstaddress := mail.Address{Name: names[0], Address: addresses[0]}
	message.AddBcc(firstaddress)
	//
	if message.Bcc[0] != firstaddress.String() {
		test.Fatal("Incorrect first element")
	}
	//
	secondaddress := mail.Address{Name: names[1], Address: addresses[1]}
	message.AddBcc(secondaddress)
	//
	if message.Bcc[1] != secondaddress.String() {
		test.Fatal("Incorrect second element")
	}
}
