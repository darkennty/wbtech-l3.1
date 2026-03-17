package notifier

import (
	"context"
	"fmt"
	"net/smtp"

	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type EmailSender struct {
	logger  zlog.Zerolog
	from    string
	to      string
	addr    string
	auth    smtp.Auth
	useSMTP bool
}

func NewEmailSender(logger zlog.Zerolog, host string, port int, user, password, recipient string) *EmailSender {
	addr := ""
	var auth smtp.Auth
	useSMTP := host != "" && port != 0 && recipient != ""
	if useSMTP {
		addr = fmt.Sprintf("%s:%d", host, port)
		if user != "" && password != "" {
			auth = smtp.PlainAuth("", user, password, host)
		}
	}
	return &EmailSender{
		logger:  logger,
		from:    user,
		to:      recipient,
		addr:    addr,
		auth:    auth,
		useSMTP: useSMTP,
	}
}

func (s *EmailSender) Send(ctx context.Context, n *model.Notification) error {
	if n.Recipient == "" {
		n.Recipient = s.to
	}

	if !s.useSMTP {
		s.logger.Info().
			Str("channel", "email").
			Str("id", n.ID).
			Str("to", n.Recipient).
			Msg("email skipped: APP_SMTP_HOST not set (stub)")
		return nil
	}

	subject := "Notification"
	body := n.Payload
	msg := "From: " + s.from + "\r\n" +
		"To: " + n.Recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body

	if err := smtp.SendMail(s.addr, s.auth, s.from, []string{n.Recipient}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}

	s.logger.Info().
		Str("channel", "email").
		Str("id", n.ID).
		Str("to", n.Recipient).
		Msg("email sent")
	return nil
}
