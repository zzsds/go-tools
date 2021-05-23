package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

const defaultStart string = "$"

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
		stream: []string{stream, defaultStart},
		count:  1,
		exit:   make(chan bool),
	}
	for _, o := range opts {
		o(sub)
	}
	if err := rdb.Ping(rdb.Context()).Err(); err != nil {
		log.Fatalf("redis ping fail %v", err)
	}
	return sub
}

// Subscribe 消费
// 目前只支持一个 stream key
// 可支持多条数据处理
func (s *subscriber) Subscribe(ctx context.Context, h event.Handler) error {
	for {
		cmd := s.reader.XRead(ctx, &redis.XReadArgs{
			Block:   s.block,
			Streams: s.stream,
			Count:   s.count,
		})
		xstream, err := cmd.Result()
		// 消息监听事件到期，自动循环
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return err
		}
		nextId := defaultStart
		// 默认只处理一个 stream ，但可以同时处理多条消息
		for _, msg := range xstream[0].Messages {
			et := event.Event{
				Key:        msg.Values["Key"].(string),
				Payload:    []byte(msg.Values["Payload"].(string)),
				Properties: map[string]string{},
			}
			json.Unmarshal([]byte(msg.Values["Properties"].(string)), &et.Properties)
			if err = h(ctx, et); err == nil {
				// 处理消息成功直接删除
				s.reader.XDel(ctx, s.stream[0], msg.ID)
			}
			nextId = msg.ID
		}
		s.stream[1] = nextId
	}
}

func (s *subscriber) Close() error {
	return s.reader.Close()
}
