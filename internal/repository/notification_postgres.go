package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type NotificationPostgresRepository struct {
	db *dbpg.DB
}

func NewNotificationPostgres(db *dbpg.DB) *NotificationPostgresRepository {
	return &NotificationPostgresRepository{db: db}
}

func (r *NotificationPostgresRepository) Create(ctx context.Context, n *model.Notification) (string, error) {
	now := time.Now().UTC()
	n.CreatedAt = now
	n.UpdatedAt = now
	n.Status = model.StatusScheduled

	const q = `
INSERT INTO notifications (
	id, channel, recipient, payload, scheduled_at, status,
	retry_count, last_error, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
`
	_, err := r.db.ExecContext(ctx, q,
		n.ID, n.Channel, n.Recipient, n.Payload, n.ScheduledAt,
		n.Status, n.RetryCount, n.LastError, n.CreatedAt, n.UpdatedAt,
	)
	return n.ID, err
}

func (r *NotificationPostgresRepository) GetNotificationByID(ctx context.Context, ID string) (*model.Notification, error) {
	const q = `
SELECT id, channel, recipient, payload, scheduled_at, status,
       retry_count, last_error, created_at, updated_at
FROM notifications
WHERE id = $1
`
	row := r.db.QueryRowContext(ctx, q, ID)
	var n model.Notification
	var lastErr sql.NullString
	if err := row.Scan(
		&n.ID, &n.Channel, &n.Recipient, &n.Payload,
		&n.ScheduledAt, &n.Status, &n.RetryCount,
		&lastErr, &n.CreatedAt, &n.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if lastErr.Valid {
		n.LastError = &lastErr.String
	}
	return &n, nil
}

func (r *NotificationPostgresRepository) GetRecentNotifications(ctx context.Context) ([]model.Notification, error) {
	const q = `
SELECT id, channel, recipient, payload, scheduled_at, status,
       retry_count, last_error, created_at, updated_at
FROM notifications
ORDER BY created_at DESC
`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print(err)
		}
	}(rows)

	var res []model.Notification
	for rows.Next() {
		var n model.Notification
		var lastErr sql.NullString
		if err = rows.Scan(
			&n.ID, &n.Channel, &n.Recipient, &n.Payload,
			&n.ScheduledAt, &n.Status, &n.RetryCount,
			&lastErr, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if lastErr.Valid {
			n.LastError = &lastErr.String
		}
		res = append(res, n)
	}
	return res, rows.Err()
}

func (r *NotificationPostgresRepository) UpdateStatus(ctx context.Context, ID string, status model.Status, retryCount int, lastErr *string) error {
	const q = `
UPDATE notifications
SET status = $2,
    retry_count = $3,
    last_error = $4,
    updated_at = $5
WHERE id = $1
`
	_, err := r.db.ExecContext(ctx, q, ID, status, retryCount, lastErr, time.Now().UTC())
	return err
}

func (r *NotificationPostgresRepository) Delete(ctx context.Context, ID string) error {
	const q = `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, ID)
	return err
}
