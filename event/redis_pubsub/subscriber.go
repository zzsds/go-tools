package redis_pubsub

import (
	"context"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

type subscriber struct {
	reader  *redis.PubSub
	channel string
}

// SubscriberOption is a subscriber option.
type SubscriberOption func(*subscriber)

// NewSubscriber new a redis subscriber.
func NewSubscriber(rdb *redis.Client, channel string) event.Subscriber {
	sub := &subscriber{
		channel: channel,
	}

	sub.reader = rdb.Subscribe(context.Background(), channel)
	return sub
}

func (s *subscriber) Subscribe(ctx context.Context, h event.Handler) error {
	for {
		msg, err := s.reader.ReceiveMessage(ctx)
		if err != nil {
			return err
		}
		_ = h(ctx, event.Event{
			Payload: []byte(msg.Payload),
		})
	}
}

func (s *subscriber) Close() error {
	return s.reader.Close()
}