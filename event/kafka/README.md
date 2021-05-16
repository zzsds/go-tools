# kafka
Kafka is a distributed event streaming platform.

## Publisher

```go
import "github.com/go-kratos/kratos/v2/event"

s := NewPublisher("test", []string{"127.0.0.1:9092"})
defer s.Close()
s.Publish(context.Background(), event.Event{Key: "key1", Payload: []byte("value1")})
s.PublishAsync(context.Background(), event.Event{Key: "key2", Payload: []byte("value2")}, nil)
s.PublishAsync(context.Background(), event.Event{Key: "key3", Payload: []byte("value3")}, func(err error) {
    log.Println(err)
})
```

## Subscriber

```go
import "github.com/go-kratos/kratos/v2/event"

r := NewSubscriber("test", "test-group", []string{"127.0.0.1:9092"})
defer r.Close()
r.Subscribe(context.Background(), func(ctx context.Context, event event.Event) error {
    log.Printf("sub: key=%s payload=%s properties=%v\n", event.Key, event.Payload, event.Properties)
    return nil
})
```
