package kafka

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/segmentio/kafka-go"
)

var (
	testTopic   = "test-event"
	testGroup   = "test-event-group"
	testBrokers = []string{"127.0.0.1:9092"}
)

func TestMain(m *testing.M) {
	// to create topics when auto.create.topics.enable='true'
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err := kafka.DialLeader(ctx, "tcp", testBrokers[0], testTopic, 0)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestPublisher(t *testing.T) {
	pub := NewPublisher(testTopic, testBrokers)
	defer pub.Close()
	pub.Publish(context.Background(), event.Event{Key: "key1", Payload: []byte("value1")})
	pub.PublishAsync(context.Background(), event.Event{Key: "key2", Payload: []byte("value2")}, nil)
	pub.PublishAsync(context.Background(), event.Event{Key: "key3", Payload: []byte("value3")}, func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})
	time.Sleep(time.Second * 2)
}
