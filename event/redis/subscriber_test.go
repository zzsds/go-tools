package redis

import (
	"context"
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
		t.Logf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
		return nil
	})
	p := NewPublisher(rdb, testChannel)
	p.Publish(context.Background(), event.Event{
		Key:     "key4",
		Payload: []byte("jayden"),
	})
}
