package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
)

func TestSubscriber(t *testing.T) {
	r := NewSubscriber(testTopic, testGroup, testBrokers)
	time.AfterFunc(time.Second*3, func() {
		r.Close()
	})
	r.Subscribe(context.Background(), func(ctx context.Context, event event.Event) error {
		t.Logf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
		return nil
	})
}
