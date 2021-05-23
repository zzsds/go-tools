package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

// ErrEventFull is a message event chan full.
var ErrEventFull = errors.New("message event chan full")

var _ event.Publisher = (*publisher)(nil)

// PublisherOption is a publisher options.
type PublisherOption func(*publisher)

func WithMaxLen(max int64) PublisherOption {
	return func(p *publisher) {
		p.maxLen = max
	}
}

func WithMaxLenApprox(max int64) PublisherOption {
	return func(p *publisher) {
		p.maxLenApprox = max
	}
}

type publisher struct {
	maxLen       int64
	maxLenApprox int64
	writer       *redis.Client
	stream       string
}

// NewPublisher new a redis publisher.
func NewPublisher(rdb *redis.Client, stream string, opts ...PublisherOption) event.Publisher {
	pub := &publisher{
		writer: rdb,
		stream: stream,
	}
	for _, o := range opts {
		o(pub)
	}
	return pub
}

func (p *publisher) Publish(ctx context.Context, e event.Event) error {
	b, _ := json.Marshal(e.Properties)
	return p.writer.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]string{
			"Key":        e.Key,
			"Payload":    string(e.Payload),
			"Properties": string(b),
		},
	}).Err()
}

func (p *publisher) Close() error {
	return p.writer.Close()
}

func (p *publisher) PublishAsync(ctx context.Context, event event.Event, callback func(err error)) error {
	return nil
}
