package smtp

import (
	"crypto/tls"

	"github.com/ztrue/tracerr"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	PasswordFile string
	SSL          bool
	From         string
}

type Sender struct {
	dialer *gomail.Dialer
	cfg    *Config
}

func NewSender(cfg *Config) (*Sender, error) {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !cfg.SSL}

	// Check connection to SMTP server
	closer, err := d.Dial()
	if err != nil {
		return nil, tracerr.Errorf("can't check smtp connection: %w", err)
	}
	if err := closer.Close(); err != nil {
		return nil, tracerr.Errorf("can't check smtp connection: %w", err)
	}

	return &Sender{
		dialer: d,
		cfg:    cfg,
	}, nil
}

func (s *Sender) SendPlainMessage(
	subject string, content string, to string, attachments ...string,
) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	return s.dialer.DialAndSend(m)
}

func (s *Sender) SendHtmlMessage(
	subject string, content string, to string, attachments ...string,
) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	return s.dialer.DialAndSend(m)
}
