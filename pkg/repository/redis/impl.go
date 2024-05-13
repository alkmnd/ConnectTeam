package redis

import (
	"ConnectTeam/pkg/repository/models"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type NotificationCache struct {
	Cache      *redis.Client
	expiration time.Duration
}

func (n NotificationCache) HGet(key string) ([]models.Notification, error) {
	var notifications []models.Notification
	ctx := context.Background()
	fields, err := n.Cache.HGetAll(ctx, key).Result()
	if err != nil {
		return notifications, err
	}
	for key, _ := range fields {
		var notification models.Notification
		err := json.Unmarshal([]byte(key), &notification)
		if err != nil {
			fmt.Printf("unmarshal error: %s\n", err)
			continue
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil

}

func (n NotificationCache) HSet(key string, value uuid.UUID, notification models.Notification) error {
	ctx := context.Background()
	err := n.Cache.HSet(ctx, key, value, notification, n.expiration).Err()
	return err
}

func (n NotificationCache) Get(key string) (models.Notification, error) {
	var notification models.Notification
	ctx := context.Background()
	err := n.Cache.Get(ctx, key).Scan(&notification)
	return notification, err
}

func NewNotificationCache(client *redis.Client) Cache {
	return &NotificationCache{
		Cache:      client,
		expiration: time.Hour * 24 * 7,
	}
}
