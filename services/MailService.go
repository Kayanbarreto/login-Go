package services

import (
	"net/smtp"
	"os"
)

type IMailService interface {
	SendEmail(email, subject, body string) error
}

type MailService struct {
}

func NewMailService() *MailService {
	return &MailService{}
}

func (service *MailService) SendEmail(email, subject, body string) error {
	// Configuração de Autenticação SMTP
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Configuração da mensagem
	message := []byte(
		"Subject: " + subject + "\r\n" +
			"\r\n" +
			body +
			"\r\n",
	)

	// Configuração de Autenticação SMTP
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Endereço de envio do email
	from := smtpUsername

	// Endereços de destino
	to := []string{email}

	// Configuração do cliente SMTP
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}
