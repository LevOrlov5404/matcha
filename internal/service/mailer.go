package service

import (
	"github.com/l-orlov/task-tracker/pkg/mailer"
	"gopkg.in/mail.v2"
)

type (
	MailerService struct {
		cfg    MailerServiceConfig
		mailer mailer.Mailer
	}
	MailerServiceConfig struct {
		From      string
		AppDomain string
	}
)

func NewMailerService(cfg MailerServiceConfig, mailer mailer.Mailer) *MailerService {
	return &MailerService{
		cfg:    cfg,
		mailer: mailer,
	}
}

func (m *MailerService) SendEmailConfirm(toEmail, token string) {
	msg := mail.NewMessage()

	msg.SetHeader("From", m.cfg.From)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Matcha registration")
	msg.SetBody("text/plain",
		"We greet you.\nTo complete the registration go by this link.\n"+
			m.cfg.AppDomain+"/confirm-email?token="+token+
			"\nThank you for choosing us :)")

	m.mailer.SendMessage(msg)
}

func (m *MailerService) SendResetPasswordConfirm(toEmail, token string) {
	msg := mail.NewMessage()

	msg.SetHeader("From", m.cfg.From)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Matcha reset password")
	msg.SetBody("text/plain",
		"Hello.\nTo reset password go by this link.\n"+
			m.cfg.AppDomain+"/confirm-reset-password?token="+token+
			"\nThank you for choosing us :)")

	m.mailer.SendMessage(msg)
}
