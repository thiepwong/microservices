package common

import (
	"bytes"
	"html/template"
	"net/smtp"
)

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
	auth    smtp.Auth
	host    string
}

func NewRequest(to []string, subject, body string, host string, sender string, auth smtp.Auth) *Request {
	return &Request{
		host:    host,
		from:    sender,
		to:      to,
		subject: subject,
		body:    body,
		auth:    auth,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := r.host

	if err := smtp.SendMail(addr, r.auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
