package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Sender interface {
	SendEmail(to, subject, body string) error
}

type LogSender struct{}

func NewLogSender() Sender {
	return &LogSender{}
}

func (s *LogSender) SendEmail(to, subject, body string) error {
	log.Printf("---------------------------------------")
	log.Printf("Sending Email to: %s", to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)
	log.Printf("---------------------------------------")
	return nil
}

type SMTPSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPSender(host string, port int, username, password, from string) Sender {
	return &SMTPSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (s *SMTPSender) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	err := smtp.SendMail(addr, auth, s.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}
	return nil
}

