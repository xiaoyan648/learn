package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"github.com/xiaoyan648/learn/delay/rabbitmq"
)

func main() {
	// # ========== 1.创建连接 ==========
	mq := rabbitmq.NewRabbitMQ()
	defer mq.Close()
	mqCh := mq.Channel

	// # ========== 2.设置队列（队列、交换机、绑定） ==========
	// 声明队列
	var err error
	_, err = mqCh.QueueDeclare(rabbitmq.Queue1, true, false, false, false, amqp.Table{
		//"x-message-ttl": 5000, // 消息过期时间（队列级别）,毫秒
	})
	rabbitmq.FailOnError(err, "创建队列失败")

	// 声明交换机
	// err = mqCh.ExchangeDeclare(Exchange1, amqp.ExchangeDirect, true, false, false, false, nil)
	err = mqCh.ExchangeDeclare(rabbitmq.Exchange1, "x-delayed-message", true, false, false, false, amqp.Table{
		"x-delayed-type": "direct",
	})
	rabbitmq.FailOnError(err, "创建交换机失败")

	// 队列绑定（将队列、routing-key、交换机三者绑定到一起）
	err = mqCh.QueueBind(rabbitmq.Queue1, rabbitmq.RoutingKey1, rabbitmq.Exchange1, false, nil)
	rabbitmq.FailOnError(err, "队列、交换机、routing-key 绑定失败")

	// # ========== 4.发布消息 ==========
	// 消息1 延迟 20s
	message := "msg" + strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(message)
	err = mqCh.Publish(rabbitmq.Exchange1, rabbitmq.RoutingKey1, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		// Expiration: "5000", // 消息过期时间（消息级别）,毫秒
		Headers: map[string]interface{}{
			"x-delay": "20000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
		},
	})
	// 消息2 延迟 5s
	message = "msg" + strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(message)
	err = mqCh.Publish(rabbitmq.Exchange1, rabbitmq.RoutingKey1, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Headers: map[string]interface{}{
			"x-delay": "5000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
		},
	})
	rabbitmq.FailOnError(err, "消息发布失败")
}
