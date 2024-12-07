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
	// 将响应转为 JSON 并发送回客户端
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Printf("Failed to send response: %v", err)
		return
	}

	log.Printf("Response sent: %s", string(responseJSON))
}

func WebSocketHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/websocket/")
		if path == "" {
			http.Error(w, "missing userId in URL", http.StatusBadRequest)
			return
		}
		userId := path // userId 直接從路徑提取
		log.Printf("Received connection for userId: %s", userId)

		// 升級到 WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}
		defer conn.Close()
		log.Printf("Connection established for userId: %s", userId)

		// 增加在線人數到 Redis
		err = redisClient.IncrementOnlineUsers()
		if err != nil {
			log.Printf("Failed to update online users: %v", err)
			return
		}
		defer redisClient.DecrementOnlineUsers()

		// 處理 WebSocket 消息
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				break
			}

			log.Printf("Raw message string: %s", string(msg))

			var message models.Message
			err = json.Unmarshal(msg, &message)
			if err != nil {
				log.Printf("Invalid message format: %v", err)
			}

			log.Printf("Parsed messages: %+v", message)

			scenario, match := service.GetScenarioByDesc(message.Mess)
			if match {
				message.Code = scenario.Code
			} else {
				message.Code = 0
			}

			log.Printf("scenario: %+v", scenario)
			if message.Code == 0 {
				response := models.Response{
					Value:      message.Mess,
					IsFinished: true,
					SubType:    "text",
					Type:       "text",
				}
				sendResponse(conn, response)
			} else {
				handler := service.GetScenarioHandler(scenario)
				options := handler.CreateOptions()
				log.Printf("Options: %+v\n", options)
				sendResponse(conn, options)
			}

			// 将消息存储到 Redis
			err = redisClient.SaveMessage(userId, message)
			if err != nil {
				log.Printf("Redis save message failed: %v", err)
			}
		}
	}
}
