1. build image
   ``` docker build -t rabbitmq .```
2. run image
   ``` docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq```
3. check
   ``` docker ps ```
4. run code
   ``` go run delay/rabbitmq/consumer/consumer.go ```
   ``` go run delay/rabbitmq/producer/producer.go ```