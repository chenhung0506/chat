package controller

import (
	"chat/internal/models"
	"chat/internal/redis"
	"chat/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func onOpen(userId string, redisClient *redis.Client) error {
	log.Printf("Connection established for userId: %s", userId)
	if err := redisClient.IncrementOnlineUsers(); err != nil {
		return err
	}
	log.Printf("Online users incremented for userId: %s", userId)

	initialMessages := models.NewInitialMessages(userId)
	if err := redisClient.SaveMessages(userId, initialMessages); err != nil {
		log.Printf("Failed to save initial message: %v", err)
		return err
	}
	log.Printf("Online users incremented for userId: %v", initialMessages)

	return nil
}

func onMessageClient(conn *websocket.Conn, clientId string, redisClient *redis.Client, message []byte) {
	messages, err := redisClient.GetMessages(clientId)
	if err != nil {
		log.Printf("Failed to get Messages: %v", err)
		return
	}

	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Invalid message format: %v", err)
		return
	}
	msg.IsClient = true

	scenario, match := service.GetScenarioByDesc(msg.Value)
	if match {
		msg.Code = scenario.Code
	} else {
		msg.Code = 0
	}

	if messages.GetPreviousMessage().Code == 6 {
		msg.Code = 6
	}

	log.Printf("previous messages: %v", messages.GetPreviousMessage())
	log.Printf("scenario: %+v", scenario)

	if messages.GetPreviousMessage().Code == 6 && !messages.GetPreviousMessage().IsFinished {
		var value = models.Message{
			UUID:       clientId,
			Value:      msg.Value,
			Code:       6,
			IsFinished: false,
			IsClient:   true,
			SubType:    "text",
			Type:       "text",
		}
		messages.AddMessage(value)
		// serviceConn, ok := models.Conn.GetService(serviceId)
		// if ok {
		// 	service.SendResponse(serviceConn, value)
		// }

		jsonBytes, err := json.Marshal(value)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}
		service.SendMessageToMQ(jsonBytes)

	} else if messages.GetPreviousMessage().Code != 0 && !messages.GetPreviousMessage().IsFinished {
		log.Println("excuteJobs===========")
		scenario = service.GetScenarioByCode(messages.GetPreviousMessage().Code)
		handler := service.GetScenarioHandler(scenario)
		value := handler.ExecuteJobs(clientId, msg.Value)
		msg.IsFinished = value.IsFinished
		messages.AddMessage(value)
		messages.AddMessage(msg)
		service.SendResponse(conn, value)
	} else {
		log.Println("CreateOptions===========")
		handler := service.GetScenarioHandler(scenario)
		value := handler.CreateOptions(clientId, redisClient)
		msg.IsFinished = value.IsFinished
		messages.AddMessage(value)
		messages.AddMessage(msg)
		service.SendResponse(conn, value)
	}

	if err := redisClient.SaveMessages(clientId, messages); err != nil {
		log.Printf("Redis save message failed: %v", err)
	}
}

func onMessageService(conn *websocket.Conn, serviceId string, redisClient *redis.Client, message []byte) {
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Invalid message format: %v", err)
		return
	}
	msg.IsClient = false

	clientId, err := redisClient.GetClientByService(serviceId)
	if err != nil {
		return
	}

	messages, err := redisClient.GetMessages(clientId)
	if err != nil {
		log.Printf("Failed to get Messages: %v", err)
		return
	}

	var value = models.Message{
		UUID:       serviceId,
		Value:      msg.Value,
		Code:       6,
		IsFinished: false,
		IsClient:   false,
		SubType:    "text",
		Type:       "text",
	}

	// clientConn, ok := models.Conn.GetClient(clientId)
	// if ok {
	// 	service.SendResponse(clientConn, value)
	// }

	messages.AddMessage(value)

	if err := redisClient.SaveMessages(clientId, messages); err != nil {
		log.Printf("Redis save message failed: %v", err)
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	service.SendMessageToMQ(jsonBytes)
}

func onClose(userId string, redisClient *redis.Client) {
	log.Printf("Connection closed for userId: %s", userId)
	if err := redisClient.DecrementOnlineUsers(); err != nil {
		log.Printf("Failed to decrement online users: %v", err)
	}
	redisClient.RemoveWaitingQueue(userId)
	log.Printf("Online users decremented for userId: %s", userId)
}

func HandleClient(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId := strings.TrimPrefix(c.Param("client_id"), "/")
		if clientId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing client_id in URL"})
			return
		}
		log.Printf("clientId: %v", clientId)

		if err := redisClient.AddClientList(clientId); err != nil {
			log.Printf("Add clientId to redis client_list fail: %v", err)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Upgrade HTTP to websocket connect error: %v ", err)
			return
		}
		defer conn.Close()

		models.Conn.AddClient(clientId, conn)

		if err := onOpen(clientId, redisClient); err != nil {
			log.Printf("Connect redis error: %v ", err)
			return
		}
		defer onClose(clientId, redisClient)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read message error: %v ", err)
				break
			}
			onMessageClient(conn, clientId, redisClient, message)
		}
	}
}

func HandleService(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceId := strings.TrimPrefix(c.Param("service_id"), "/")
		if serviceId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing service_id in URL"})
			return
		}
		log.Printf("userId: %v", serviceId)

		if err := redisClient.AddServiceList(serviceId); err != nil {
			log.Printf("Add userId to redis client_list fail: %v", err)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Upgrade HTTP to websocket connect error: %v ", err)
			return
		}
		defer conn.Close()

		models.Conn.AddService(serviceId, conn)

		if err := onOpen(serviceId, redisClient); err != nil {
			log.Printf("Connect redis error: %v ", err)
			return
		}
		defer onClose(serviceId, redisClient)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read message error: %v ", err)
				break
			}
			onMessageService(conn, serviceId, redisClient, message)
		}
	}
}
