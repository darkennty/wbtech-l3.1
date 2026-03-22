package service

import (
	"context"

	"WBTech_L3.1/internal/cache"
	"WBTech_L3.1/internal/config"
	"WBTech_L3.1/internal/model"
	"WBTech_L3.1/internal/notifier"
	"WBTech_L3.1/internal/queue"
	"WBTech_L3.1/internal/repository"
)

type Notification interface {
	Create(ctx context.Context, req CreateRequest) (string, error)
	GetNotificationByID(ctx context.Context, ID string) (*model.Notification, error)
	GetRecentNotifications(ctx context.Context) ([]model.Notification, error)
	UpdateStatus(ctx context.Context, ID string, status model.Status, retryCount int, lastErr *string) error
	Delete(ctx context.Context, ID string) error
	CancelNotification(ctx context.Context, id string) error
	HandleMessage(ctx context.Context, msg queue.Message) error
}

type Service struct {
	Notification
}

func NewService(repo *repository.Repository, queue queue.Queue, senders notifier.Senders, cache cache.StatusCache, cfg config.Config) *Service {
	return &Service{
		Notification: NewNotificationService(repo.Notification, queue, senders, cache, cfg),
	}
}
