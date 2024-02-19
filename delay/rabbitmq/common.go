package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	Queue1      = "queue1"
	Exchange1   = "exchange1"
	RoutingKey1 = "key1"
	MqUrl       = "amqp://guest:guest@localhost:5672/"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// 拿到 rabbitmq 的 channel，轻量级 connection
func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial(MqUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}
}

func (r RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}

// FailOnError 错误处理函数
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
