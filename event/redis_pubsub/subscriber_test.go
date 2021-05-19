package redis_pubsub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
)

func TestSubscriber(t *testing.T) {
	r := NewSubscriber(rdb, testChannel)
	time.AfterFunc(time.Second*3, func() {
		r.Close()
	})
	go r.Subscribe(context.Background(), func(ctx context.Context, event event.Event) error {
		fmt.Println(event, string(event.Payload), "33333")
		// t.Logf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
		return nil
	})
	time.Sleep(2 * time.Second)
	p := NewPublisher(rdb, testChannel)
	p.Publish(context.Background(), event.Event{
		// Key:     testChannel,
		Payload: []byte(`{"amount": 12, "id": 1}`),
	})
}
