package redis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

type subscriber struct {
	reader *redis.Client
	stream []string
	block  time.Duration
	count  int64
	exit   chan bool
}

// SubscriberOption is a subscriber option.
type SubscriberOption func(*subscriber)

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
// 将错误重新存到List 中
func (s *subscriber) Subscribe(ctx context.Context, h event.Handler) error {
	for {
		cmd := s.reader.XRead(ctx, &redis.XReadArgs{
			Block:   s.block,
			Streams: s.stream,
			Count:   s.count,
		})
		stream, err := cmd.Result()
		if err != nil {
			return err
		}
		var xsm redis.XStream
		if len(stream) > 0 {
			xsm = stream[0]
		}
		var msg redis.XMessage
		if len(xsm.Messages) > 0 {
			msg = xsm.Messages[0]
		}
		// next stream
		s.stream[1] = msg.ID
		et := event.Event{}
		ekey := xsm.Stream + "Error"
		mb, _ := json.Marshal(msg)
		switch val := msg.Values[xsm.Stream].(type) {
		case string:
			b, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				log.Printf("ID %s base64 decode fail %s listErr %v", msg.ID, err.Error(), s.reader.LPush(ctx, ekey, mb).Err())
				continue
			}
			if err := json.Unmarshal(b, &et); err != nil {
				log.Printf("ID %s Unmarshal fail %s listErr %v", msg.ID, err.Error(), s.reader.LPush(ctx, ekey, mb).Err())
				continue
			}
		default:
			// 直接跳出
			continue
		}
		if err := h(ctx, et); err != nil {
			log.Printf("ID %s Handle fail %s listErr %v", msg.ID, err.Error(), s.reader.LPush(ctx, ekey, mb).Err())
		}
	}
}

func (s *subscriber) Close() error {
	return s.reader.Close()
}
