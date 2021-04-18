package service

import (
	"crypto/tls"

	"github.com/l-orlov/matcha/internal/config"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

type (
	MailerService struct {
		cfg           config.Mailer
		log           *logrus.Entry
		dialer        *gomail.Dialer
		msgToSendChan chan *gomail.Message
	}
)

func NewMailerService(cfg config.Mailer, log *logrus.Entry) *MailerService {
	d := gomail.NewDialer(
		cfg.ServerAddress.Host, cfg.ServerAddress.Port, cfg.Username, cfg.Password,
	)
	d.Timeout = cfg.Timeout.Duration()
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.ServerAddress.Host,
		InsecureSkipVerify: false,
	}

	mailerSvc := &MailerService{
		cfg:    cfg,
		log:    log,
		dialer: d,
	}

	mailerSvc.msgToSendChan = make(chan *gomail.Message, cfg.MsgToSendChanSize)
	mailerSvc.InitWorkers()

	return mailerSvc
}

func (s *MailerService) InitWorkers() {
	for i := 0; i < s.cfg.WorkersNum; i++ {
		go func() {
			for m := range s.msgToSendChan {
				if err := s.dialer.DialAndSend(m); err != nil {
					s.log.Errorf("failed to send message by email: %v", err)
				}
			}
		}()
	}
}

func (s *MailerService) Close() {
	close(s.msgToSendChan)
}

func (s *MailerService) SendEmailConfirm(toEmail, token string) {
	m := gomail.NewMessage()

	m.SetHeader("From", s.cfg.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Matcha registration")
	m.SetBody("text/plain",
		"We greet you.\nTo complete the registration go by this link.\n"+
			"localhost:8080/confirm-email?token="+token+
			"\nThank you for choosing us :)")

	s.msgToSendChan <- m
}

func (s *MailerService) SendResetPasswordConfirm(toEmail, token string) {
	m := gomail.NewMessage()

	m.SetHeader("From", s.cfg.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Matcha reset password")
	m.SetBody("text/plain",
		"Hello.\nTo reset password go by this link.\n"+
			"localhost:8080/confirm-reset-password?token="+token+
			"\nThank you for choosing us :)")

	s.msgToSendChan <- m
}
