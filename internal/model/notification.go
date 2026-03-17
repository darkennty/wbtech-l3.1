package model

import "time"

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type Notification struct {
	ID          string    `json:"id"`
	Channel     string    `json:"channel"`
	Recipient   string    `json:"recipient"`
	Payload     string    `json:"payload"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      Status    `json:"status"`
	RetryCount  int       `json:"retry_count"`
	LastError   *string   `json:"last_error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
