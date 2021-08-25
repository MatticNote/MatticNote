package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"net/smtp"
	"strings"
)

func SendMail(to, subject, contentType, bodyStr string) error {
	auth := smtp.PlainAuth(
		"",
		config.Config.Mail.Username,
		config.Config.Mail.Password,
		config.Config.Mail.SmtpHost,
	)

	var (
		header, body, message bytes.Buffer
	)

	header.WriteString(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n",
		config.Config.Mail.From,
		to,
	))
	header.WriteString(encodeSubject(subject))
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", contentType))
	header.WriteString("Content-Transfer-Encoding: base64\r\n")

	body.WriteString(bodyStr)

	message = header
	message.WriteString("\r\n")
	message.WriteString(add76crlf(base64.StdEncoding.EncodeToString(body.Bytes())))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Config.Mail.SmtpHost, config.Config.Mail.SmtpPort),
		auth,
		config.Config.Mail.From,
		[]string{to},
		[]byte(message.String()),
	)
}

func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

func utf8Split(utf8string string, length int) []string {
	var resultString []string
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}
