//
// Package email allows to send emails with attachments.
//
//	With base in scorredoira project
//
//	(c) 2020 2ww Engenharia de Sistemas
//
//	wilson.souza
//
package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

//	define email body
const (
	EmailTextType = "text/plain"
	EmailHTMLType = "text/html"
)

// Attachment represents an email attachment.
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// Header represents an additional email header.
type Header struct {
	Key   string
	Value string
}

// Message represents a smtp message.
type Message struct {
	From            mail.Address
	To              []string
	Cc              []string
	Bcc             []string
	ReplyTo         string
	Subject         string
	Body            string
	BodyContentType string
	Headers         []Header
	Attachments     map[string]*Attachment
}

// Attach attaches a file.
func (m *Message) Attach(file string, inline bool) error {
	data, failure := ioutil.ReadFile(file)
	if failure != nil {
		return failure
	}
	//
	_, filename := filepath.Split(file)
	//
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}
	return nil
}

func (m *Message) AddTo(address mail.Address) []string {
	m.To = append(m.To, address.String())
	return m.To
}

func (m *Message) AddCc(address mail.Address) []string {
	m.Cc = append(m.Cc, address.String())
	return m.Cc
}

func (m *Message) AddBcc(address mail.Address) []string {
	m.Bcc = append(m.Bcc, address.String())
	return m.Bcc
}

// AttachBuffer attaches a binary attachment.
func (m *Message) AttachBuffer(filename string, buf []byte, inline bool) error {
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

// Inline includes a file as an inline attachment.
func (m *Message) Inline(file string) error {
	return m.Attach(file, true)
}

// Ads a Header to message
func (m *Message) AddHeader(key string, value string) Header {
	newHeader := Header{Key: key, Value: value}
	m.Headers = append(m.Headers, newHeader)
	return newHeader
}

func CreateMessage(subject, body, contenttype string) *Message {
	message := &Message{Subject: subject, Body: body, BodyContentType: contenttype}
	message.Attachments = make(map[string]*Attachment)
	return message
}

// Tolist returns all the recipients of the email
func (m *Message) ToList() []string {
	rcptList := []string{}

	toList, _ := mail.ParseAddressList(strings.Join(m.To, ","))
	for _, to := range toList {
		rcptList = append(rcptList, to.Address)
	}

	ccList, _ := mail.ParseAddressList(strings.Join(m.Cc, ","))
	for _, cc := range ccList {
		rcptList = append(rcptList, cc.Address)
	}

	bccList, _ := mail.ParseAddressList(strings.Join(m.Bcc, ","))
	for _, bcc := range bccList {
		rcptList = append(rcptList, bcc.Address)
	}

	return rcptList
}

// Bytes returns the mail data in bytes array
func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.From.String() + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	// Add custom headers
	if len(m.Headers) > 0 {
		for _, header := range m.Headers {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", header.Key, header.Value))
		}
	}

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
					buf.WriteString(mime)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

// Send sends the message.
func Send(address string, authentication smtp.Auth, message *Message) error {
	return smtp.SendMail(address,
		authentication,
		message.From.Address,
		message.ToList(),
		message.ToBytes())
}
