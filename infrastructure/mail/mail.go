package mail

import (
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/knadh/smtppool/v2"
)

func NewMailSender(settings entities.Settings) (Sender, error) {
	smtpSettings := settings.SMTPSettings

	opts := smtppool.Opt{
		Host:            smtpSettings.Host,
		Port:            smtpSettings.Port,
		MaxConns:        smtpSettings.MaxConnections,
		IdleTimeout:     smtpSettings.IdleTimeout,
		PoolWaitTimeout: smtpSettings.PoolWaitTimeout,
		SSL:             smtppool.SSLNone,
	}

	pool, err := smtppool.New(opts)
	if err != nil {
		return nil, derr.JoinError("failed to create SMTP pool", err)
	}

	return mailSender{
		pool:     pool,
		settings: smtpSettings,
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
