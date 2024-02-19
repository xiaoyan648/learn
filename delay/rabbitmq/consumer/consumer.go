package main

import (
	"log"
	"time"

	"github.com/xiaoyan648/learn/delay/rabbitmq"
)

func main() {
	// # ========== 1.创建连接 ==========
	mq := rabbitmq.NewRabbitMQ()
	defer mq.Close()
	mqCh := mq.Channel

	// # ========== 2.消费消息 ==========
	msgsCh, err := mqCh.Consume(rabbitmq.Queue1, "", false, false, false, false, nil)
	rabbitmq.FailOnError(err, "消费队列失败")

	forever := make(chan bool)
	go func() {
		for d := range msgsCh {
			// 要实现的逻辑
			log.Printf("接收的消息: %s, time: %s", d.Body, time.Now().Local().Format("2006-01-02 15:04:05"))

			// 手动应答
			d.Ack(false)
			// d.Reject(true)
		}
	}()
	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
}
