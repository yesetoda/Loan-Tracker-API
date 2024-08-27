package infrastructures

import (
	"log"
	"net/smtp"
)

func SendEmail(toEmail string, title string, body string, link string) error {

	message := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>` + title + `</title>
	</head>
	<body>
		<h1>` + title + `</h1>
		<p>` + body + `</p>
		<a href="` + link + `">Click the Link</a>
	</body>
	</html>
	`

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	key := "kdjb ggie qdei gqhl"
	host := "smtp.gmail.com"
	auth := smtp.PlainAuth("", "yeneineh.seiba@a2sv.org", key, "smtp.gmail.com")

	port := "587"
	address := host + ":" + port
	messages := []byte(mime + message)

	err := smtp.SendMail(address, auth, "abel.wendmu@a2sv.org", []string{toEmail}, messages)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
