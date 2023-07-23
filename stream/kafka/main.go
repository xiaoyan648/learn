package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"api.version.request":           "true",
		"message.max.bytes":             1000000,
		"linger.ms":                     500,
		"sticky.partitioning.linger.ms": 1000,
		"retries":                       2147483647 - 1000,
		"retry.backoff.ms":              1000,
		"acks":                          "1",

		"bootstrap.servers":                   "47.101.195.2:9093,47.101.35.245:9093,101.132.135.104:9093",
		"security.protocol":                   "sasl_ssl",
		"ssl.ca.location":                     "/Users/litao/code/learn/stream/kafka/conf/kafka-cert",
		"sasl.username":                       "bigdatashanghaidev",
		"sasl.password":                       "bi6GshAg8Nha9dEv",
		"sasl.mechanism":                      "PLAIN",
		"enable.ssl.certificate.verification": false,
	})
	if err != nil {
		panic(err)
	}
	go func() {
		topic := "new-media-event"
		for {
			c := make(chan kafka.Event, 1)
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny,
				},
				Value: []byte("this is demo"),
			}, c)
			fmt.Println(err)
			event := <-c
			fmt.Println(event)
			time.Sleep(time.Second)
		}
	}()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":         "47.101.195.2:9093,47.101.35.245:9093,101.132.135.104:9093",
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000,
		"group.id":                  "new-media-event-test-group",
		//"auto.offset.reset":                   "earliest",
		"security.protocol":                   "SASL_SSL",
		"ssl.ca.location":                     "/Users/litao/code/learn/stream/kafka/conf/kafka-cert",
		"sasl.username":                       "bigdatashanghaidev",
		"sasl.password":                       "bi6GshAg8Nha9dEv",
		"sasl.mechanism":                      "PLAIN",
		"enable.ssl.certificate.verification": false,
	})
	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"new-media-event"}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		fmt.Println("start read message")
		// msg, err := c.ReadMessage(3 * time.Second)
		ev := c.Poll(3000 * int(time.Microsecond.Microseconds()))
		// if err == nil {
		// 	fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		// } else {
		// 	// The client will automatically try to recover from all errors.
		// 	// Timeout is not considered an error because it is raised by
		// 	// ReadMessage in absence of messages.
		// 	fmt.Printf("Consumer error: %v (%v)\n", err.(kafka.Error).IsTimeout(), msg)
		// }
		fmt.Printf("end read message %+v", ev)
	}

	c.Close()
}
