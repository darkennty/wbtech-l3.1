package notifier

import (
	"context"
	"fmt"

	"WBTech_L3.1/internal/config"
	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type Sender interface {
	Send(ctx context.Context, n *model.Notification) error
}

type Senders map[string]Sender

func NewSenders(cfg config.Config, logger zlog.Zerolog) Senders {
	s := make(Senders)
	s["email"] = NewEmailSender(logger, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.EmailDefaultRecipient)
	s["telegram"] = NewTelegramSender(logger, cfg.TelegramBotToken, cfg.TelegramDefaultRecipient)
	return s
}

func (s Senders) Get(channel string) (Sender, error) {
	sender, ok := s[channel]
	if !ok {
		return nil, fmt.Errorf("unsupported channel: %s", channel)
	}
	return sender, nil
}
