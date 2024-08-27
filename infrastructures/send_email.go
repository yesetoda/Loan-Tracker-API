package infrastructures

import (
	"example/b/Loan-Tracker-API/config"
	"fmt"
	"log"
	"net/smtp"
)

// SendEmail sends an email with the specified title, body, and link to the specified email address.
func SendEmail(toEmail, title, body, link string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return err
	}

	// Construct the HTML message
	message := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
		</head>
		<body>
			<h1>%s</h1>
			<p>%s</p>
			<a href="%s">Click the Link</a>
		</body>
		</html>
	`, title, title, body, link)

	// Prepare the MIME header and message body
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	emailMessage := []byte(mime + message)

	// Setup SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	address := smtpHost + ":" + smtpPort

	auth := smtp.PlainAuth("", cfg.Email.EmailSender, cfg.Email.EmailKey, smtpHost)

	// Send the email
	err = smtp.SendMail(address, auth, cfg.Email.EmailSender, []string{toEmail}, emailMessage)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
		return err
	}

	log.Printf("Email successfully sent to %s", toEmail)
	return nil
}
