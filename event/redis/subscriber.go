package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

type subscriber struct {
	reader *redis.Client
	stream []string
	group  string
	block  time.Duration
	count  int64
	exit   chan bool
}

// SubscriberOption is a subscriber option.
type SubscriberOption func(*subscriber)

// WithGroup 组内消费
func WithGroup(group string) SubscriberOption {
	return func(s *subscriber) {
		s.group = group
	}
}

// WithBlock 阻塞时长
// block 0 表示永远阻塞，直到消息到来，block 1000 表示阻塞 1s，如果 1s 内没有任何消息到来，就返回 nil
func WithBlock(block time.Duration) SubscriberOption {
	return func(s *subscriber) {
		s.block = block
	}
}

// WithCount 查询条数
func WithCount(count int64) SubscriberOption {
	return func(s *subscriber) {
		s.count = count
	}
}

// WithStream
func WithStream(stream ...string) SubscriberOption {
	return func(s *subscriber) {
		s.stream = stream
	}
}

// NewSubscriber new a redis subscriber.
func NewSubscriber(rdb *redis.Client, stream string, opts ...SubscriberOption) event.Subscriber {
	sub := &subscriber{
		reader: rdb,
		block:  5 * time.Second,
		stream: []string{stream, "$"},
		count:  1,
		exit:   make(chan bool),
	}
	for _, o := range opts {
		o(sub)
	}
	return sub
}

// Subscribe 消费
func (s *subscriber) Subscribe(ctx context.Context, h event.Handler) error {
	go func() {
		for {
			select {
			case <-s.exit:
				return
			default:
				cmd := s.reader.XRead(ctx, &redis.XReadArgs{
					Block:   s.block,
					Streams: s.stream,
					Count:   s.count,
				})
				if cmd.Err() != nil {
					fmt.Println(cmd.Err())
					continue
				}
				s, err := cmd.Result()
				if err != nil {
					log.Fatalln(err)
					continue
				}
				fmt.Println(s)
				if err := h(ctx, event.Event{Key: "jayden"}); err != nil {
					fmt.Println(err)
				}
				// s.reader.XAck(ctx, s.stream[0], "", cmd.)
			}
		}
	}()
	return nil
}

func (s *subscriber) Close() error {
	s.exit <- true
	close(s.exit)
	return s.reader.Close()
}
