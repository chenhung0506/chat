package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"chat/internal/models"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr, password string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &Client{rdb: rdb}
}

func (c *Client) IncrementOnlineUsers() error {
	return c.rdb.Incr(context.Background(), "online_number").Err()
}

func (c *Client) DecrementOnlineUsers() error {
	return c.rdb.Decr(context.Background(), "online_number").Err()
}

func (c *Client) AddWaitingQueue(clientId string) {
	_, err := c.rdb.SAdd(context.Background(), "waiting_queue", clientId).Result()
	if err != nil {
		log.Printf("err: %v", err)
	}
}

func (c *Client) GetWaitingQueue() ([]string, error) {
	data, err := c.rdb.SMembers(context.Background(), "waiting_queue").Result()
	if err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}
	return data, nil
}

func (c *Client) RemoveWaitingQueue(clientId string) {
	log.Printf("remove clientId: %v from waiting_queue", clientId)
	_, err := c.rdb.SRem(context.Background(), "waiting_queue", clientId).Result()
	if err != nil {
		log.Printf("err: %v", err)
	}
}

func (c *Client) RemoveWholeWaitingQueue() error {
	_, err := c.rdb.Del(context.Background(), "waiting_queue").Result()
	return err
}

func (c *Client) GetMessages(userId string) (models.Messages, error) {
	key := fmt.Sprintf("messages:%s", userId)
	data, err := c.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return models.Messages{}, err
	}

	var messages models.Messages
	err = json.Unmarshal([]byte(data), &messages)
	if err != nil {
		return models.Messages{}, err
	}
	return messages, nil
}

func (c *Client) SaveMessages(userId string, messages models.Messages) error {
	key := fmt.Sprintf("messages:%s", userId)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return c.rdb.Set(context.Background(), key, data, 0).Err()
}

func (c *Client) SaveClientAndService(cliendId string, serviceId string) error {
	err1 := c.rdb.HSet(context.Background(), "client_to_service", cliendId, serviceId).Err()
	if err1 != nil {
		return err1
	}
	err2 := c.rdb.HSet(context.Background(), "service_to_client", serviceId, cliendId).Err()
	if err2 != nil {
		return err2
	}
	return nil
}

func (c *Client) GetClientByService(serviceId string) (string, error) {
	clientId, err := c.rdb.HGet(context.Background(), "service_to_client", serviceId).Result()
	if err != nil {
		log.Printf("GetClientByService error: %v", err)
		return "", err
	}
	return clientId, nil
}

func (c *Client) GetServiceByClient(clientId string) (string, error) {
	serviceId, err := c.rdb.HGet(context.Background(), "client_to_service", clientId).Result()
	if err != nil {
		log.Printf("GetServiceByClient error: %v", err)
		return "", err
	}
	return serviceId, nil
}

func (c *Client) AddClientList(cliendId string) error {
	return c.rdb.SAdd(context.Background(), "client_list", cliendId).Err()
}

func (c *Client) AddServiceList(serviceId string) error {
	return c.rdb.SAdd(context.Background(), "service_list", serviceId).Err()
}

func (c *Client) GetClientList() []string {
	users, _ := c.rdb.SMembers(context.Background(), "client_list").Result()
	return users
}

func (c *Client) GetServiceList() []string {
	services, _ := c.rdb.SMembers(context.Background(), "service_list").Result()
	return services
}

// 客戶從人工客服清單移除
func (c *Client) RemoveFromClientList(clientId string) {
	c.rdb.SRem(context.Background(), "client_list", clientId)
}
