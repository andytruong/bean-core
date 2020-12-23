package model

import (
	"gopkg.in/gomail.v2"
)

type Message struct {
	From      string `json:"from"`
	Recipient string `json:"recipient"`
	Body      Body   `json:"body"`
}

type Body struct {
	Subject     string   `json:"subject"`
	HTML        string   `json:"html"`
	Attachments []string `json:"attachments"`
}

func (m Message) Send(dialer *gomail.Dialer) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.Recipient)
	msg.SetHeader("Subject", m.Body.Subject)
	msg.SetBody("text/html", m.Body.HTML)

	for _, path := range m.Body.Attachments {
		msg.Attach(path)
	}

	return dialer.DialAndSend(msg)
}
