package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

// Email struct holds the configuration for sending emails
type Email struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

// NewEmailSender initializes and returns a new Email service
func NewEmailSender(host string, port string, username string, password string, sender string) *Email {
	return &Email{
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUsername: username,
		SMTPPassword: password,
		SMTPFrom:     sender,
	}
}

// SendEmail sends an email using the provided template and payload
func (e *Email) SendEmail(payload interface{}, pathTemplate string, to []string, subject string) error {
	// Get the absolute path of the template file
	absPath, err := filepath.Abs(pathTemplate)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	// Parse the template file
	t, err := template.ParseFiles(absPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with the payload
	buff := new(bytes.Buffer)
	if err := t.Execute(buff, payload); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// Prepare the email message
	msg := "To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		buff.String()

	fmt.Println("test", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)
	// Set up authentication
	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	// Send the email
	if err := smtp.SendMail(e.SMTPHost+":"+e.SMTPPort, auth, e.SMTPFrom, to, []byte(msg)); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
