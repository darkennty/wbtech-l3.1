package repository

import (
	"context"

	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type Notification interface {
	Create(ctx context.Context, notification *model.Notification) (string, error)
	GetNotificationByID(ctx context.Context, ID string) (*model.Notification, error)
	GetRecentNotifications(ctx context.Context, limit int) ([]model.Notification, error)
	UpdateStatus(ctx context.Context, ID string, status model.Status, retryCount int, lastErr *string) error
	Delete(ctx context.Context, ID string) error
}

type Repository struct {
	Notification
}

func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{
		Notification: NewNotificationPostgres(db),
	}
}
