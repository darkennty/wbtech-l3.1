package service

import (
	"context"
	"math"
	"time"

	"WBTech_L3.1/internal/cache"
	"WBTech_L3.1/internal/config"
	"WBTech_L3.1/internal/model"
	"WBTech_L3.1/internal/notifier"
	"WBTech_L3.1/internal/queue"
	"WBTech_L3.1/internal/repository"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo    repository.Notification
	queue   queue.Queue
	senders notifier.Senders
	cache   cache.StatusCache
	cfg     config.Config
}

type CreateRequest struct {
	Channel   string
	Recipient string
	Message   string
	SendAt    time.Time
}

func NewNotificationService(repo repository.Notification, queue queue.Queue, senders notifier.Senders, cache cache.StatusCache, cfg config.Config) *NotificationService {
	return &NotificationService{
		repo:    repo,
		queue:   queue,
		senders: senders,
		cache:   cache,
		cfg:     cfg,
	}
}

func (s *NotificationService) Create(ctx context.Context, req CreateRequest) (string, error) {
	recipient := req.Recipient
	if recipient == "" {
		if req.Channel == "telegram" {
			recipient = s.cfg.TelegramDefaultRecipient
		} else if req.Channel == "email" {
			recipient = s.cfg.EmailDefaultRecipient
		}
	}

	n := &model.Notification{
		ID:          uuid.NewString(),
		Channel:     req.Channel,
		Recipient:   recipient,
		Payload:     req.Message,
		ScheduledAt: req.SendAt,
		Status:      model.StatusScheduled,
	}
	if _, err := s.repo.Create(ctx, n); err != nil {
		return "", err
	}
	_ = s.cache.Set(ctx, n)

	msg := queue.Message{
		NotificationID: n.ID,
		ExecuteAt:      req.SendAt,
		RetryCount:     0,
	}
	if err := s.queue.Publish(ctx, msg); err != nil {
		return "", err
	}
	return n.ID, nil
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, ID string) (*model.Notification, error) {
	if n, err := s.cache.Get(ctx, ID); err == nil && n != nil {
		return n, nil
	}
	n, err := s.repo.GetNotificationByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, n)
	return n, nil
}

func (s *NotificationService) GetRecentNotifications(ctx context.Context) ([]model.Notification, error) {
	return s.repo.GetRecentNotifications(ctx)
}

func (s *NotificationService) UpdateStatus(ctx context.Context, ID string, status model.Status, retryCount int, lastErr *string) error {
	return s.repo.UpdateStatus(ctx, ID, status, retryCount, lastErr)
}

func (s *NotificationService) Delete(ctx context.Context, ID string) error {
	if err := s.repo.Delete(ctx, ID); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, ID)

	return nil
}

func (s *NotificationService) CancelNotification(ctx context.Context, id string) error {
	n, err := s.repo.GetNotificationByID(ctx, id)
	if err != nil {
		return err
	}
	if n.Status == model.StatusSent || n.Status == model.StatusFailed || n.Status == model.StatusCancelled {
		if err = s.repo.Delete(ctx, id); err != nil {
			return err
		}
		_ = s.cache.Delete(ctx, id)
		return nil
	}

	if err = s.repo.UpdateStatus(ctx, id, model.StatusCancelled, n.RetryCount, n.LastError); err != nil {
		return err
	}
	n.Status = model.StatusCancelled
	_ = s.cache.Set(ctx, n)

	return nil
}

func (s *NotificationService) HandleMessage(ctx context.Context, msg queue.Message) error {
	n, err := s.repo.GetNotificationByID(ctx, msg.NotificationID)
	if err != nil {
		return err
	}
	if n.Status == model.StatusCancelled || n.Status == model.StatusSent {
		return nil
	}
	now := time.Now().UTC()
	if now.Before(msg.ExecuteAt) {
		delay := msg.ExecuteAt.Sub(now)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		n, err = s.repo.GetNotificationByID(ctx, msg.NotificationID)
		if err != nil {
			return err
		}
		if n.Status == model.StatusCancelled || n.Status == model.StatusSent {
			return nil
		}
	}

	sender, err := s.senders.Get(n.Channel)
	if err != nil {
		return s.fail(ctx, n, msg, err)
	}

	if err = sender.Send(ctx, n); err != nil {
		return s.fail(ctx, n, msg, err)
	}

	if err = s.repo.UpdateStatus(ctx, n.ID, model.StatusSent, msg.RetryCount, nil); err != nil {
		return err
	}
	n.Status = model.StatusSent
	n.RetryCount = msg.RetryCount
	n.LastError = nil
	_ = s.cache.Set(ctx, n)

	return nil
}

func (s *NotificationService) fail(ctx context.Context, n *model.Notification, msg queue.Message, sendErr error) error {
	msg.RetryCount++
	lastErr := sendErr.Error()

	if msg.RetryCount > s.cfg.MaxRetryCount {
		if err := s.repo.UpdateStatus(ctx, n.ID, model.StatusFailed, msg.RetryCount, &lastErr); err != nil {
			return err
		}
		n.Status = model.StatusFailed
		n.RetryCount = msg.RetryCount
		n.LastError = &lastErr
		_ = s.cache.Set(ctx, n)
		return nil
	}

	delay := s.retryDelay(msg.RetryCount)
	executeAt := time.Now().UTC().Add(delay)
	msg.ExecuteAt = executeAt

	if err := s.repo.UpdateStatus(ctx, n.ID, model.StatusPending, msg.RetryCount, &lastErr); err != nil {
		return err
	}
	n.Status = model.StatusPending
	n.RetryCount = msg.RetryCount
	n.LastError = &lastErr
	_ = s.cache.Set(ctx, n)

	return s.queue.Publish(ctx, msg)
}

func (s *NotificationService) retryDelay(retryCount int) time.Duration {
	if retryCount <= 0 {
		return s.cfg.BaseRetryDelay
	}
	mul := math.Pow(2, float64(retryCount-1))
	return time.Duration(float64(s.cfg.BaseRetryDelay) * mul)
}
