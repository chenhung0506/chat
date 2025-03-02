package service

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendResponse(conn *websocket.Conn, response interface{}) {
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
