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

func (this Message) Send(dialer *gomail.Dialer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", this.From)
	m.SetHeader("To", this.Recipient)
	m.SetHeader("Subject", this.Body.Subject)
	m.SetBody("text/html", this.Body.HTML)

	for _, path := range this.Body.Attachments {
		m.Attach(path)
	}

	return dialer.DialAndSend(m)
}
