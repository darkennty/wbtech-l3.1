package config

import (
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	HTTPAddr                 string
	DatabaseDSN              string
	RabbitURL                string
	RabbitQueue              string
	RabbitExchange           string
	RedisAddr                string
	RedisPassword            string
	RedisDB                  int
	BaseRetryDelay           time.Duration
	MaxRetryCount            int
	EmailFrom                string
	SMTPHost                 string
	SMTPPort                 int
	SMTPUser                 string
	SMTPPassword             string
	EmailDefaultRecipient    string
	TelegramBotToken         string
	TelegramDefaultRecipient string
}

func Load() Config {
	c := config.New()

	_ = c.LoadEnvFiles(".env")
	c.EnableEnv("APP")

	c.SetDefault("http.addr", "8080")
	c.SetDefault("db.dsn", "postgres://postgres:postgres@localhost:5432/delayed_notifier?sslmode=disable")
	c.SetDefault("rabbit.url", "amqp://guest:guest@localhost:5672/")
	c.SetDefault("rabbit.queue", "notifications")
	c.SetDefault("rabbit.exchange", "notifications")
	c.SetDefault("redis.addr", "localhost:6379")
	c.SetDefault("redis.password", "")
	c.SetDefault("redis.db", 0)
	c.SetDefault("retry.base_delay", time.Minute)
	c.SetDefault("retry.max", 5)
	c.SetDefault("email.from", "no-reply@example.com")
	c.SetDefault("smtp.host", "")
	c.SetDefault("smtp.port", 587)
	c.SetDefault("smtp.user", "")
	c.SetDefault("smtp.password", "")
	c.SetDefault("telegram.token", "")
	c.SetDefault("default.recipient", "")

	return Config{
		HTTPAddr:                 c.GetString("http.addr"),
		DatabaseDSN:              c.GetString("db.dsn"),
		RabbitURL:                c.GetString("rabbit.url"),
		RabbitQueue:              c.GetString("rabbit.queue"),
		RabbitExchange:           c.GetString("rabbit.exchange"),
		RedisAddr:                c.GetString("redis.addr"),
		RedisPassword:            c.GetString("redis.password"),
		RedisDB:                  c.GetInt("redis.db"),
		BaseRetryDelay:           c.GetDuration("retry.base_delay"),
		MaxRetryCount:            c.GetInt("retry.max"),
		EmailFrom:                c.GetString("email.from"),
		SMTPHost:                 c.GetString("smtp.host"),
		SMTPPort:                 c.GetInt("smtp.port"),
		SMTPUser:                 c.GetString("smtp.user"),
		SMTPPassword:             c.GetString("smtp.password"),
		EmailDefaultRecipient:    c.GetString("email.default.recipient"),
		TelegramBotToken:         c.GetString("telegram.token"),
		TelegramDefaultRecipient: c.GetString("telegram.default.recipient"),
	}
}
