package redis

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
)

func TestSubscriber(t *testing.T) {
	r := NewSubscriber(rdb, testChannel, WithBlock(-1))
	time.AfterFunc(time.Second*4, func() {
		r.Close()
	})
	go func() {
		err := r.Subscribe(context.Background(), func(ctx context.Context, event event.Event) error {
			// fmt.Println(event, string(event.Payload), "33333")
			t.Logf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)
	p := NewPublisher(rdb, testChannel)
	p.Publish(context.Background(), event.Event{
		Key:        testChannel,
		Payload:    []byte(`{"amount": 12, "id": 1}`),
		Properties: map[string]string{"name": "jayden"},
	})

	time.Sleep(2 * time.Second)
}

func TestXInfoGroups(t *testing.T) {
	cmd := rdb.XInfoGroups(context.Background(), testChannel)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}
	t.Log(cmd.Result())
}
