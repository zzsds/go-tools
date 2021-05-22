package redis

import (
	"context"
	"encoding/base64"
	"reflect"
	"time"
	"unsafe"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

var defaultID string = "$"

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
		stream: []string{stream, defaultID},
		count:  1,
		exit:   make(chan bool),
	}
	for _, o := range opts {
		o(sub)
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
		if err != nil {
			return err
		}
		nextId := defaultID
		// 同时处理多条消息
		for _, msg := range xstream[0].Messages {
			val, _ := msg.Values[xstream[0].Stream].(string)
			b, _ := base64.StdEncoding.DecodeString(val)
			_ = h(ctx, *(*event.Event)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)))
		}
		s.stream[1] = nextId
	}
}

func (s *subscriber) Close() error {
	return s.reader.Close()
}
