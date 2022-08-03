package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zzsds/go-tools/event"
)

var (
	_ event.Sender   = (*sender)(nil)
	_ event.Receiver = (*receiver)(nil)
)

// SenderOption is a sender options.
type SenderOption func(*sender)

func WithMaxLen(max int64) SenderOption {
	return func(p *sender) {
		p.maxLen = max
	}
}

func WithMaxLenApprox(max int64) SenderOption {
	return func(p *sender) {
		p.maxLenApprox = max
	}
}

type sender struct {
	maxLen       int64
	maxLenApprox int64
	writer       *redis.Client
	stream       string
}

// NewSender new a redis sender.
func NewSender(rdb *redis.Client, stream string, opts ...SenderOption) event.Sender {
	pub := &sender{
		writer: rdb,
		stream: stream,
	}
	for _, o := range opts {
		o(pub)
	}
	return pub
}

func (p *sender) Send(ctx context.Context, e event.Event) error {
	return p.writer.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]interface{}{
			"Key":     e.Key,
			"Payload": string(e.Value()),
		},
	}).Err()
}

func (p *sender) Close() error {
	return p.writer.Close()
}

const defaultStart string = "$"

type receiver struct {
	reader *redis.Client
	stream []string
	block  time.Duration
	count  int64
	exit   chan bool
}

// ReceiverOption is a receiver option.
type ReceiverOption func(*receiver)

// WithBlock 阻塞时长
// block 0 表示永远阻塞，直到消息到来，block 1000 表示阻塞 1s，如果 1s 内没有任何消息到来，就返回 nil
func WithBlock(block time.Duration) ReceiverOption {
	return func(s *receiver) {
		s.block = block
	}
}

// WithCount 查询条数
func WithCount(count int64) ReceiverOption {
	return func(s *receiver) {
		s.count = count
	}
}

// WithStream
func WithStream(stream ...string) ReceiverOption {
	return func(s *receiver) {
		s.stream = stream
	}
}

// NewReceiver new a redis receiver.
func NewReceiver(rdb *redis.Client, stream string, opts ...ReceiverOption) event.Receiver {
	sub := &receiver{
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

// Receive 接收消息
// 目前只支持一个 stream key
// 可支持多条数据处理
func (s *receiver) Receive(ctx context.Context, handler event.Handler) error {
	go func() {
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
				break
			}
			nextId := defaultStart
			// 默认只处理一个 stream ，但可以同时处理多条消息
			for _, msg := range xstream[0].Messages {
				err = handler(ctx, event.NewMessage(msg.Values["Key"].(string), []byte(msg.Values["Payload"].(string))))
				if err != nil {
					log.Fatal("message handling exception:", err)
				}

				// 处理消息成功直接删除
				if err := s.reader.XDel(ctx, s.stream[0], msg.ID).Err(); err != nil {
					log.Fatal("failed to del messages:", err)
				}
				nextId = msg.ID
			}
			s.stream[1] = nextId
		}
	}()
	return nil
}

func (s *receiver) Close() error {
	return s.reader.Close()
}
