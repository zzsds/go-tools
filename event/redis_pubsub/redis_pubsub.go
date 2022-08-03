package redis_pubsub

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/zzsds/go-tools/event"
)

var (
	_ event.Sender   = (*sender)(nil)
	_ event.Receiver = (*receiver)(nil)
)

// SenderOption is a sender options.
type SenderOption func(*sender)

type sender struct {
	writer  *redis.Client
	channel string
}

// NewSender new a redis sender.
func NewSender(rdb *redis.Client, channel string) event.Sender {
	pub := &sender{rdb, channel}
	return pub
}

func (p *sender) Send(ctx context.Context, event event.Event) error {
	return p.writer.Publish(ctx, p.channel, redis.Message{
		Channel: p.channel,
		Pattern: event.Key(),
		Payload: string(event.Value()),
	}).Err()
}

func (p *sender) Close() error {
	return p.writer.Close()
}

type receiver struct {
	reader  *redis.PubSub
	channel string
}

// ReceiverOption is a receiver option.
type ReceiverOption func(*receiver)

// NewReceiver new a redis receiver.
func NewReceiver(rdb *redis.Client, channel string) event.Receiver {
	sub := &receiver{
		channel: channel,
	}

	sub.reader = rdb.Subscribe(context.Background(), channel)
	return sub
}

func (s *receiver) Receive(ctx context.Context, handler event.Handler) error {
	go func() {
		for {
			m, err := s.reader.ReceiveMessage(ctx)
			if err != nil {
				break
			}

			err = handler(context.Background(), event.NewMessage(string(m.Pattern), []byte(m.Payload)))
			if err != nil {
				log.Fatal("message handling exception:", err)
			}
		}
	}()
	return nil
}

func (s *receiver) Close() error {
	return s.reader.Close()
}
