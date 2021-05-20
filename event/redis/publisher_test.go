package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/event"
	"github.com/go-redis/redis/v8"
)

var (
	testPass    = "123456"
	testChannel = "AssetFrozen"
	testArr     = "127.0.0.1:6379"
	rdb         *redis.Client
)

func TestMain(m *testing.M) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     testArr,
		Password: testPass,
		DB:       0,
	})
	os.Exit(m.Run())
}

func TestPublisher(t *testing.T) {
	pub := NewPublisher(rdb, testChannel)
	defer pub.Close()
	if err := pub.Publish(context.Background(), event.Event{Key: testChannel, Payload: []byte(`{"id": 1, "amount": 12}`)}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
}

func TestPublishers(t *testing.T) {
	pub := NewPublisher(rdb, testChannel)
	// defer pub.Close()
	for i := 1; i < 102; i++ {
		payload := fmt.Sprintf(`{"id": %d, "amount": %d}`, i, i*10)
		if err := pub.Publish(context.Background(), event.Event{Key: testChannel, Payload: []byte(payload)}); err != nil {
			t.Fatal(err)
		}
		// t.Log(payload)
	}
}
