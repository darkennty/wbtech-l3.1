package worker

import (
	"context"

	"WBTech_L3.1/internal/queue"
	"WBTech_L3.1/internal/service"

	"github.com/rs/zerolog"
)

func Run(ctx context.Context, svc *service.Service, q queue.Queue, logger zerolog.Logger) {
	if err := q.StartConsuming(ctx, func(hctx context.Context, msg queue.Message) error {
		logger.Info().
			Str("notification_id", msg.NotificationID).
			Time("execute_at", msg.ExecuteAt).
			Int("retry_count", msg.RetryCount).
			Msg("worker: processing message")

		if err := svc.HandleMessage(hctx, msg); err != nil {
			logger.Error().Err(err).Str("notification_id", msg.NotificationID).Msg("failed to handle message")
			return err
		}

		logger.Info().Str("notification_id", msg.NotificationID).Msg("worker: notification sent")
		return nil
	}); err != nil {
		if ctx.Err() != nil {
			logger.Info().Msg("worker stopped")
			return
		}
		logger.Error().Err(err).Msg("failed to start consuming queue")
	}
}
