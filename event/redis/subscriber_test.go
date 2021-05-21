package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
)

func TestSubscribers(t *testing.T) {
	r := NewSubscriber(rdb, testChannel)
	go r.Subscribe(context.Background(), nil)
	t.Run("TestPublishers", TestPublishers)
	time.Sleep(5 * time.Second)
}

func TestSubscriber(t *testing.T) {
	r := NewSubscriber(rdb, testChannel)
	err := r.Subscribe(context.Background(), func(ctx context.Context, event event.Event) error {
		fmt.Println(event, string(event.Payload), "33333")
		t.Logf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestXInfoStream(t *testing.T) {
	cmd := rdb.XInfoStream(context.Background(), testChannel)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}
	t.Log(cmd.Result())
}

func TestXInfoGroups(t *testing.T) {
	cmd := rdb.XInfoGroups(context.Background(), testChannel)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}
	t.Log(cmd.Result())
}
