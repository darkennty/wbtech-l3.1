package repository

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

func NewPostgresDB(ctx context.Context, dsn string) (*dbpg.DB, error) {
	db, err := dbpg.New(dsn, nil, &dbpg.Options{
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = db.Master.PingContext(ctx); err != nil {
		_ = db.Master.Close()
		for _, s := range db.Slaves {
			_ = s.Close()
		}
		return nil, err
	}

	return db, nil
}
