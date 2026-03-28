package util

import (
	"tessera/config"
	"tessera/models"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendVerificationEmail(user *models.User, data *EmailData) error {
	config, err := config.LoadConfig("./")

	if err != nil {
		return fmt.Errorf("could not load config %v", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return fmt.Errorf("could not parse template %v", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// func SendVerificationEmail(user *models.User, data *EmailData) error {
// 	config, err := config.LoadConfig("./")
// 	if err != nil {
// 		return fmt.Errorf("could not load config %v", err)
// 	}

// 	// Create a buffered channel to send emails concurrently
// 	emailChan := make(chan *gomail.Message, 1)

// 	// Create a buffered channel to send errors from goroutines to main function
// 	errChan := make(chan error, 1)

// 	// Create a single SMTP dialer instance
// 	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	// Generate email content
// 	var body bytes.Buffer
// 	template, err := ParseTemplateDir("templates")
// 	if err != nil {
// 		return fmt.Errorf("could not parse template %v", err)
// 	}
// 	template.ExecuteTemplate(&body, "verificationCode.html", &data)

// 	// Create an email message
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", config.EmailFrom)
// 	m.SetHeader("To", user.Email)
// 	m.SetHeader("Subject", data.Subject)
// 	m.SetBody("text/html", body.String())
// 	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

// 	// Add the email message to the channel
// 	emailChan <- m

// 	go func() {
// 		// if err := d.DialAndSend(<-emailChan); err != nil {
// 		// 	fmt.Println("Failed to send email:", err)
// 		// }
// 		if err := d.DialAndSend(<-emailChan); err != nil {
// 			errChan <- fmt.Errorf("failed to send email: %v", err)
// 			return
// 		}
// 		errChan <- nil
// 	}()

// 	// Check if there is an error
// 	if err := <-errChan; err != nil {
// 		return err
// 	}

// 	return nil
// }
