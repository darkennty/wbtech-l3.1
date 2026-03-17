package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/zlog"
)

const telegramAPI = "https://api.telegram.org"

type TelegramSender struct {
	logger zlog.Zerolog
	token  string
	userID string
	client *http.Client
}

func NewTelegramSender(logger zlog.Zerolog, token, userID string) *TelegramSender {
	return &TelegramSender{
		logger: logger,
		token:  token,
		userID: userID,
		client: &http.Client{},
	}
}

func (s *TelegramSender) Send(ctx context.Context, n *model.Notification) error {
	if n.Recipient == "" {
		n.Recipient = s.userID
	}

	if s.token == "" || n.Recipient == "" {
		s.logger.Info().
			Str("channel", "telegram").
			Str("id", n.ID).
			Str("to", n.Recipient).
			Msg("telegram skipped: APP_TELEGRAM_TOKEN or APP_TELEGRAM_DEFAULT_RECIPIENT not set")
		return nil
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPI, s.token)
	body := map[string]string{
		"chat_id": n.Recipient,
		"text":    n.Payload,
	}
	raw, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api: status %d", resp.StatusCode)
	}

	var out struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if !out.OK {
		return fmt.Errorf("telegram api: %s", out.Description)
	}

	s.logger.Info().
		Str("channel", "telegram").
		Str("id", n.ID).
		Str("to", n.Recipient).
		Msg("telegram sent")
	return nil
}
