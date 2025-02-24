package controller

import (
	"chat/internal/redis"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWaitingClients(c *gin.Context, redisClient *redis.Client) {
	waitingClients, err := redisClient.GetWaitingQueue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得等待中的客戶"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"waiting_clients": waitingClients})
}

func RemoveWholeWaitingQueue(c *gin.Context, redisClient *redis.Client) {
	err := redisClient.RemoveWholeWaitingQueue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Clear waiting_queue fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clear waiting_queue success"})
}

func AssignClient(c *gin.Context, redisClient *redis.Client) {
	var request struct {
		ClientID  string `json:"client_id"`
		ServiceID string `json:"service_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供正確的 client_id 與 service_id"})
		return
	}

	waitingClients, err := redisClient.GetWaitingQueue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得等待中的客戶"})
		return
	}
	found := false
	for _, id := range waitingClients {
		if id == request.ClientID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此客戶不在等待列表中"})
		return
	}

	err = redisClient.SaveClientAndService(request.ClientID, request.ServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得等待中的客戶"})
		return
	}

	log.Printf("Service %s assigned to Client %s", request.ServiceID, request.ClientID)

	messages, err := redisClient.GetMessages(request.ClientID)
	if err != nil {
		log.Printf("Failed to get Messages: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功指派客服", "client_id": request.ClientID, "service_id": request.ServiceID, "data": messages})
}
