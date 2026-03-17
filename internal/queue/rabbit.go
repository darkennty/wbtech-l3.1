package queue

import (
	"context"
	"encoding/json"
	"time"

	"WBTech_L3.1/internal/cache"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type Message struct {
	NotificationID string    `json:"notification_id"`
	ExecuteAt      time.Time `json:"execute_at"`
	RetryCount     int       `json:"retry_count"`
}

type Rabbit struct {
	Client *rabbitmq.RabbitClient
}

func NewRabbit(url string) (*Rabbit, error) {
	client, err := rabbitmq.NewClient(rabbitmq.ClientConfig{
		URL: url,
		ReconnectStrat: retry.Strategy{
			Attempts: 0,
			Delay:    1 * time.Second,
			Backoff:  2,
		},
		ProducingStrat: retry.Strategy{Attempts: 3, Delay: 200 * time.Millisecond, Backoff: 2},
		ConsumingStrat: retry.Strategy{Attempts: 3, Delay: 200 * time.Millisecond, Backoff: 2},
	})
	if err != nil {
		return nil, err
	}
	return &Rabbit{Client: client}, nil
}

type Queue interface {
	Publish(ctx context.Context, msg Message) error
	Cancel(ctx context.Context, notificationID string) error
	StartConsuming(ctx context.Context, handler func(context.Context, Message) error) error
}

type rabbitQueue struct {
	client    *rabbitmq.RabbitClient
	queueName string
	exchange  string
	pub       *rabbitmq.Publisher
	canceller cache.ICancellationService
}

func NewRabbitQueue(client *rabbitmq.RabbitClient, canceller cache.ICancellationService, exchange, queueName string) Queue {
	return &rabbitQueue{
		client:    client,
		queueName: queueName,
		exchange:  exchange,
		pub:       rabbitmq.NewPublisher(client, exchange, "application/json"),
		canceller: canceller,
	}
}

func (q *rabbitQueue) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return q.pub.Publish(ctx, body, q.queueName)
}

func (q *rabbitQueue) Cancel(ctx context.Context, notificationID string) error {
	return q.canceller.Cancel(ctx, notificationID)
}

func (q *rabbitQueue) StartConsuming(ctx context.Context, handler func(context.Context, Message) error) error {
	cons := rabbitmq.NewConsumer(q.client, rabbitmq.ConsumerConfig{
		Queue:         q.queueName,
		ConsumerTag:   "delayed-notifier",
		AutoAck:       false,
		Workers:       4,
		PrefetchCount: 10,
		Nack: rabbitmq.NackConfig{
			Multiple: false,
			Requeue:  true,
		},
		Ask: rabbitmq.AskConfig{Multiple: false},
	}, func(ctx context.Context, d amqp091.Delivery) error {
		var msg Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			return nil
		}

		isCancelled, err := q.canceller.IsCancelled(ctx, msg.NotificationID)
		if err != nil {
			return err
		}

		if isCancelled {
			return nil
		}

		err = handler(ctx, msg)

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return nil
			}
			return err
		}

		return nil
	})
	return cons.Start(ctx)
}
