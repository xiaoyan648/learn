package main

import (
	"context"
	"fmt"
	"time"

	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// 1.rocketmq
func TestRocketMQ() {
	// 1. 部署 https://github.com/apache/rocketmq/tree/master
	// 2. 代码 https://github.com/apache/rocketmq-client-go/blob/master/examples/consumer/delay/main.go
	close, err := RunRocketConsumer()
	if err != nil {
		panic(err)
	}
	defer close()

	if err := SendDelayWithRocket(); err != nil {
		panic(err)
	}

	time.Sleep(30 * time.Second)
}

func SendDelayWithRocket() error {
	p, err := rocketmq.NewProducer(
		producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		// producer.WithNsResolver(primitive.NewPassthroughResolver(endPoint)),
		producer.WithRetry(2),
		producer.WithGroupName("test_delay"),
	)
	if err != nil {
		return err
	}

	if err := p.Start(); err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		msg := primitive.NewMessage("delay", []byte("this is a delay message!"))
		msg.WithDelayTimeLevel(2)
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}

	return p.Shutdown()
}

func RunRocketConsumer() (close func() error, err error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("test_delay"),
		consumer.WithNsResovler(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	)
	if err != nil {
		return func() error { return nil }, err
	}

	err = c.Subscribe("delay", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt,
	) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			t := time.Now().UnixNano()/int64(time.Millisecond) - msg.BornTimestamp
			fmt.Printf("Receive message[msgId=%s] %d ms later\n", msg.MsgId, t)
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return func() error { return nil }, err
	}

	if err := c.Start(); err != nil {
		return func() error { return nil }, err
	}

	return c.Shutdown, nil
}
