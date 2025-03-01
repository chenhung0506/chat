package controller

import (
	"chat/internal/models"
	"chat/internal/redis"
	"chat/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendResponse(conn *websocket.Conn, response interface{}) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
	log.Printf("Response sent: %s", string(responseJSON))
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

	serviceId, err := redisClient.GetServiceByClient(clientId)
	if err == nil {
		log.Printf("seand message from clientId: %v to serviceId: %v", clientId, serviceId)
		serviceConn, exists := services[serviceId]
		var value = models.Message{
			Value:      msg.Value,
			Code:       6,
			IsFinished: true,
			IsClient:   msg.IsClient,
			SubType:    "text",
			Type:       "text",
		}
		if exists {
			sendResponse(serviceConn, value)
			messages.AddMessage(msg)
		} else {
			// conn, ch := Rabbitmq.connectRabbitMQ()
			// defer conn.Close()
			// defer ch.Close()
		}

		return
	}

	scenario, match := service.GetScenarioByDesc(msg.Value)
	if match {
		msg.Code = scenario.Code
	} else {
		msg.Code = 0
	}

	log.Printf("previous messages: %v", messages.GetPreviousMessage())
	log.Printf("scenario: %+v", scenario)

	if messages.GetPreviousMessage().Code != 0 && !messages.GetPreviousMessage().IsFinished {
		log.Println("excuteJobs===========")
		scenario = service.GetScenarioByCode(messages.GetPreviousMessage().Code)
		handler := service.GetScenarioHandler(scenario)
		value := handler.ExecuteJobs(msg.Value)
		msg.IsFinished = value.IsFinished
		messages.AddMessage(value)
		messages.AddMessage(msg)
		sendResponse(conn, value)
	} else {
		log.Println("CreateOptions===========")
		handler := service.GetScenarioHandler(scenario)
		options := handler.CreateOptions(clientId, redisClient)
		msg.IsFinished = options.IsFinished
		messages.AddMessage(options)
		messages.AddMessage(msg)
		sendResponse(conn, options)
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

	var clientConn = clients[clientId]

	var value = models.Message{
		Value:      msg.Value,
		Code:       6,
		IsFinished: true,
		IsClient:   msg.IsClient,
		SubType:    "text",
		Type:       "text",
	}
	sendResponse(clientConn, value)
	messages.AddMessage(msg)

	if err := redisClient.SaveMessages(clientId, messages); err != nil {
		log.Printf("Redis save message failed: %v", err)
	}
}

func onClose(userId string, redisClient *redis.Client) {
	log.Printf("Connection closed for userId: %s", userId)
	if err := redisClient.DecrementOnlineUsers(); err != nil {
		log.Printf("Failed to decrement online users: %v", err)
	}
	redisClient.RemoveWaitingQueue(userId)
	log.Printf("Online users decremented for userId: %s", userId)
}

var (
	clients    = make(map[string]*websocket.Conn)
	services   = make(map[string]*websocket.Conn)
	clientsMu  sync.Mutex
	servicesMu sync.Mutex
)

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

		clientsMu.Lock()
		clients[clientId] = conn
		clientsMu.Unlock()

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

		servicesMu.Lock()
		services[serviceId] = conn
		servicesMu.Unlock()
		service_list := redisClient.GetServiceList()
		log.Printf("client_list: %v", service_list)

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
