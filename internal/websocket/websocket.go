package websocket

import (
	"chat/internal/models"
	"chat/internal/redis"
	"chat/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

func onMessage(conn *websocket.Conn, userId string, redisClient *redis.Client, message []byte) {
	messages, err := redisClient.GetMessages(userId)
	if err != nil {
		log.Printf("Failed to get Messages: %v", err)
		return
	}

	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Invalid message format: %v", err)
		return
	}

	scenario, match := service.GetScenarioByDesc(msg.Mess)
	if match {
		msg.Code = scenario.Code
	} else {
		msg.Code = 0
	}

	log.Printf("previous messages: %v", messages.GetPreviousMessage())
	log.Printf("scenario: %+v", scenario)

	if messages.GetPreviousMessage().Code != 0 && !messages.GetPreviousMessage().IsFinished {
		log.Println("excuteJobs:::::::::::::")
		scenario = service.GetScenarioByCode(messages.GetPreviousMessage().Code)
		handler := service.GetScenarioHandler(scenario)
		value := handler.ExecuteJobs(msg.Mess)
		msg.IsFinished = value.IsFinished
		sendResponse(conn, value)
	} else {
		log.Println("CreateOptions:::::::::::::")
		handler := service.GetScenarioHandler(scenario)
		options := handler.CreateOptions()
		msg.IsFinished = options.IsFinished
		sendResponse(conn, options)
	}
	log.Println("===============================")
	messages.AddMessage(msg)

	if err := redisClient.SaveMessages(userId, messages); err != nil {
		log.Printf("Redis save message failed: %v", err)
	}
}

func onClose(userId string, redisClient *redis.Client) {
	log.Printf("Connection closed for userId: %s", userId)
	if err := redisClient.DecrementOnlineUsers(); err != nil {
		log.Printf("Failed to decrement online users: %v", err)
	}
	log.Printf("Online users decremented for userId: %s", userId)
}

func onError(err error) {
	log.Printf("WebSocket error: %v", err)
}

func WebSocketHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/websocket/")
		if path == "" {
			http.Error(w, "missing userId in URL", http.StatusBadRequest)
			return
		}
		userId := path

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			onError(err)
			return
		}
		defer conn.Close()

		if err := onOpen(userId, redisClient); err != nil {
			onError(err)
			return
		}
		defer onClose(userId, redisClient)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				onError(err)
				break
			}
			onMessage(conn, userId, redisClient, message)
		}
	}
}
