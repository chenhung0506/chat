package service

import (
	"chat/internal/models"
	"chat/internal/redis"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// 連接到 RabbitMQ
func connectRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@139.162.2.175:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return conn, ch
}

// 創建 Fanout Exchange
func SetupExchange() {
	conn, ch := connectRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	// err := ch.ExchangeDeclare(
	// 	"chat_exchange", "fanout",
	// 	true, false, false, false, nil,
	// )

	err := ch.ExchangeDeclare(
		"chat_fanout_exchange", // Exchange 名稱
		"fanout",               // 類型：Fanout
		true,                   // durable 設為 true（持久化）
		false,                  // auto-delete
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	log.Println("Exchange 'chat_fanout_exchange' created")
}

func SendMessageToMQ(message []byte) {
	conn, ch := connectRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	err := ch.Publish(
		"chat_fanout_exchange", // 交換機名稱
		"",                     // routing key（Fanout 不使用 routing key）
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message sent: %s", message)
}

func ReceiveMessage(queueName string, redisClient *redis.Client) {
	conn, ch := connectRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	// queue, err := ch.QueueDeclare(
	// "",
	// false,
	// true,
	// true,
	// false,
	// nil,
	// )

	_, err := ch.QueueDeclare(
		queueName, // 隊列名稱
		true,      // durable 設為 true（持久化）
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// err = ch.QueueBind(
	// queue.Name,
	// "",
	// "chat_fanout_exchange",
	// false,
	// nil)

	err = ch.QueueBind(
		queueName,              // 隊列名稱
		"",                     // routing key（Fanout 不需要）
		"chat_fanout_exchange", // 交換機名稱
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // 監聽自己的隊列
		"",        // consumer name
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Printf("%s is waiting for messages...", queueName)
	for msg := range msgs {
		log.Printf("[%s] Received message: %s", queueName, msg.Body)

		var message models.Message
		err = json.Unmarshal([]byte(msg.Body), &message)
		if err != nil {
			log.Printf("Unmarshal rabbitMQ message error: %v", err)
		}

		if message.IsClient {
			serviceId, err := redisClient.GetServiceByClient(message.UUID)
			if err == nil {
				log.Printf("seand message from clientId: %v to serviceId: %v", message.UUID, serviceId)
				serviceConn, ok := models.Conn.GetService(serviceId)
				if ok {
					SendResponse(serviceConn, message)
				}
			}
		} else {
			clientId, err := redisClient.GetClientByService(message.UUID)
			if err == nil {
				log.Printf("seand message from serviceId: %v to clientId: %v", message.UUID, clientId)
				clientConn, ok := models.Conn.GetClient(clientId)
				if ok {
					SendResponse(clientConn, message)
				}
			}
		}

	}
}
