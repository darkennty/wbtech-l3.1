package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"WBTech_L3.1/internal/api/handler"
	"WBTech_L3.1/internal/api/server"
	"WBTech_L3.1/internal/cache"
	"WBTech_L3.1/internal/config"
	"WBTech_L3.1/internal/notifier"
	"WBTech_L3.1/internal/queue"
	"WBTech_L3.1/internal/repository"
	"WBTech_L3.1/internal/service"
	"WBTech_L3.1/internal/worker"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	_ = os.Setenv("TZ", "UTC")

	zlog.InitConsole()
	logger := zlog.Logger
	cfg := config.Load()

	db, err := repository.NewPostgresDB(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}
	defer func() {
		_ = db.Master.Close()
		for _, s := range db.Slaves {
			_ = s.Close()
		}
	}()

	if err = repository.Migrate(context.Background(), db); err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate database")
	}

	rabbit, err := queue.NewRabbit(cfg.RabbitURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init rabbitmq")
	}
	defer func() {
		_ = rabbit.Client.Close()
	}()

	if err = rabbit.Client.DeclareExchange(cfg.RabbitExchange, "direct", true, false, false, nil); err != nil {
		logger.Fatal().Err(err).Msg("failed to declare rabbit exchange")
	}
	if err = rabbit.Client.DeclareQueue(cfg.RabbitQueue, cfg.RabbitExchange, cfg.RabbitQueue, true, false, true, nil); err != nil {
		logger.Fatal().Err(err).Msg("failed to declare rabbit queue")
	}

	redisClient := cache.NewRedisClient(cfg)
	defer func() {
		if redisClient != nil {
			_ = redisClient.Close()
		}
	}()

	canceller := cache.NewRedisCancellationService(redisClient)
	senders := notifier.NewSenders(cfg, logger)
	q := queue.NewRabbitQueue(rabbit.Client, canceller, cfg.RabbitExchange, cfg.RabbitQueue)

	statusCache := cache.NewStatusCache(redisClient)
	repo := repository.NewRepository(db)
	services := service.NewService(repo, q, senders, statusCache, cfg)

	srv := new(server.Server)
	handlers := handler.NewHandler(services, logger)

	go func() {
		logger.Info().Str("addr", cfg.HTTPAddr).Msg("starting http server")
		if err = srv.Run(cfg.HTTPAddr, handlers.InitRoutes()); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatal().Err(err).Msg("http server error")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	workerCtx, workerCancel := context.WithCancel(context.Background())
	go worker.Run(workerCtx, services, q, logger)

	<-ctx.Done()
	logger.Info().Msg("shutting down")
	workerCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http server shutdown error")
	}
}
