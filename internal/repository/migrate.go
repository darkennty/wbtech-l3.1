package repository

import (
	"context"

	"github.com/wb-go/wbf/dbpg"
)

func Migrate(ctx context.Context, db *dbpg.DB) error {
	const q = `
CREATE TABLE IF NOT EXISTS notifications (
	id UUID PRIMARY KEY,
	channel TEXT NOT NULL,
	recipient TEXT NOT NULL,
	payload TEXT NOT NULL,
	scheduled_at TIMESTAMPTZ NOT NULL,
	status TEXT NOT NULL,
	retry_count INT NOT NULL DEFAULT 0,
	last_error TEXT,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);
`
	_, err := db.ExecContext(ctx, q)
	return err
}
