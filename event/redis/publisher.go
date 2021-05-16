package redis

import (
	"bytes"
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

// ErrEventFull is a message event chan full.
var ErrEventFull = errors.New("message event chan full")

var _ event.Publisher = (*publisher)(nil)

// PublisherOption is a publisher options.
type PublisherOption func(*publisher)

type publisher struct {
	writer  *redis.Client
	channel string
}

// NewPublisher new a redis publisher.
func NewPublisher(rdb *redis.Client, channel string) event.Publisher {
	pub := &publisher{rdb, channel}
	return pub
}

func (p *publisher) Publish(ctx context.Context, event event.Event) error {
	return p.writer.Publish(ctx, p.channel, bytes.NewBuffer(event.Payload).String()).Err()
}

func (p *publisher) Close() error {
	return p.writer.Close()
}

func (p *publisher) PublishAsync(ctx context.Context, event event.Event, callback func(err error)) error {
	return nil
}
