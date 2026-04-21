package mail

import (
	"net/smtp"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/knadh/smtppool"
)

func NewSender(settings entities.SMTPSettings) (Sender, error) {
	opts := smtppool.Opt{
		Host:            settings.Host,
		Port:            settings.Port,
		MaxConns:        settings.MaxConnections,
		IdleTimeout:     settings.IdleTimeout,
		PoolWaitTimeout: settings.PoolWaitTimeout,
		Auth:            smtp.PlainAuth("", settings.From, settings.Password, settings.Host),
	}

	pool, err := smtppool.New(opts)
	if err != nil {
		return nil, derr.JoinError("failed to create SMTP pool", err)
	}

	return mailSender{
		pool:     pool,
		settings: settings,
	}, nil
}

type mailSender struct {
	pool     *smtppool.Pool
	settings entities.SMTPSettings
}

func (s mailSender) SendMail(to []string, subject, content string) error {
	email := smtppool.Email{
		From:    s.settings.From,
		To:      to,
		Subject: subject,
		Text:    []byte(content),
	}

	err := s.pool.Send(email)
	if err != nil {
		return derr.JoinError("failed to send email", err)
	}

	return nil
}
