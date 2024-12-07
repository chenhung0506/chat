package redis

import (
	"context"
	"encoding/json"
	"fmt"

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
	return c.rdb.Incr(context.Background(), "OnlineNumber").Err()
}

func (c *Client) DecrementOnlineUsers() error {
	return c.rdb.Decr(context.Background(), "OnlineNumber").Err()
}

func (c *Client) SaveMessage(uuid string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("user:%s:messages", uuid)
	return c.rdb.RPush(context.Background(), key, data).Err()
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
